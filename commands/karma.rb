require 'redis'

class Karma
  def initialize(nickname = nil)
    @nickname = nickname
    unmute
  end

  def rem
    karma_was = muted { score }

    count = redis.zrem redis_key, normalized_nickname
    say "```#{@nickname} has just been removed, (s)he had #{karma_was.to_i} points of karma```"

    count
  end

  def incr(by = 1)
    karma = redis.zincrby redis_key, by, normalized_nickname
    say "```New karma for #{@nickname} is #{karma.to_i}```"

    karma
  end

  def decr(by = 1)
    incr(-by)
  end

  def score
    karma = redis.zscore redis_key, normalized_nickname
    say "```#{@nickname} has #{karma.to_i} points of karma```"

    karma
  end

  def rank
    position = redis.zrevrank redis_key, normalized_nickname

    if position.nil?
      say "```#{@nickname} is not in the list!```"
      return
    end

    say "```#{@nickname} is at #{ordinalize(position + 1)} position```"

    position
  end

  def range(start = 0, stop = -1)
    ranking = redis.zrevrange redis_key, start, stop, with_scores: true

    max_nickname_width = ranking.map { |x| x.first.length }.max
    max_score_width    = ranking.map { |x| x.last.to_i.to_s.length }.max

    say "```"
    ranking.each do |nickname, score|
      say [
        nickname.ljust(max_nickname_width, "."),
        score.to_i.to_s.rjust(max_score_width, " ")
      ].join(" ")
    end
    say "```"

    ranking
  end

  def winners
    range(0, 4)
  end

  def losers
    range(-5, -1)
  end

  def identify(user_id)
    redis.hset("#{redis_key}:idbyname", user_id, @nickname)
  end

  def fetch_nickname
    nickname = redis.hget("#{redis_key}:idbyname", @nickname)
    unless nickname.nil?
      @nickname = nickname
    end
  end

  private

  def say(s)
    puts s unless muted?
  end

  def muted?
    @muted
  end

  def mute
    @muted = true
  end

  def unmute
    @muted = false
  end

  def muted
    muted_was = @muted

    @muted = true
    res = yield

    @muted = muted_was

    res
  end

  def redis
    @redis ||= Redis.new(db: 7)
  end

  def redis_key
    "cbot:karma"
  end

  def normalized_nickname
    @nickname.downcase.sub(/^@/, "")
  end

  def ordinalize(position)
    ord = if (11..13).include? position
            "th"
          else
            { 1 => "st", 2 => "nd", 3 => "rd" }.fetch(position % 10, "th")
          end

    "#{position}#{ord}"
  end
end