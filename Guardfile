def run(cmdline)
  puts "+#{cmdline}"
  system cmdline
end

guard :shell do
  watch /\.go$/ do |m|
    puts "#{Time.now}: #{m[0]}"
    case m[0]
    when /_test\.go$/
      parent = File.dirname m[0]
      sources = Dir["#{parent}/*.go"].reject{|p| %w(_test.go _windows.go).any?{|s| p.end_with? s } }.join(' ')
      run "go test -v #{m[0]} #{sources}"
      run "golint #{m[0]}"
    else
      run 'go build ./'
      run "golint #{m[0]}"
    end
  end
end
