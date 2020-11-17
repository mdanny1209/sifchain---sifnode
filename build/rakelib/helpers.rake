#
# Current working directory.
#
def cwd
  File.dirname(__FILE__)
end

#
# Check the supplied arguments
#
# @param args Arguments passed to rake
#
def check_args(args)
  if args[:chainnet] == nil
    puts "Please provider a chainnet argument E.g testnet, mainnet, etc"
    exit
  end

  case args[:provider]
  when "aws"
  when "az"
    puts "Build me!"
    exit
  when "gcp"
    puts "Build me!"
    exit
  when "do"
    puts "Build me!"
    exit
  else
    puts "Please provide a cloud host provider. E.g aws"
    exit
  end
end

#
# Network config
#
# @params chainnet Name or ID of the chain
#
def network_config(chainnet)
  "#{cwd}/../networks/#{Digest::SHA256.hexdigest chainnet}.yml"
end

#
# Generic prompt
#
# @param args Arguments passed to rake
#
def are_you_sure(args)
  if args[:skip_prompt].nil?
    STDOUT.puts "Are you sure? (y/n)"

    begin
      input = STDIN.gets.strip.downcase
    end until %w(y n).include?(input)

    exit(0) if input != 'y'
  end
end
