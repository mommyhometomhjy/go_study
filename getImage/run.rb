require 'roo'


xlsx = Roo::Spreadsheet.open(File.dirname(__FILE__)+'/资料.xlsx')
sheet = xlsx.sheet(0)

sheet.each(no: "货号", title: "标题", key: "关键字", price: "价格", material: "材质", image: "图片", note: "备注") do |hash|
  next if hash[:no] == "货号"

  puts "采集#{hash[:no]}图片"
  Dir.mkdir(File.dirname(__FILE__)+"/#{hash[:no]}") unless Dir.exist?(File.dirname(__FILE__)+"/#{hash[:no]}")


  # puts ship_sku[0].class
  # puts ship_sku
 
  #生产txt
  File.open(File.dirname(__FILE__)+"/#{hash[:no]}/#{hash[:no]}.txt", 'w') do |f|
    f.puts hash[:title].gsub(hash[:no], "M" + hash[:no]) if hash[:title]
    f.puts "M" + hash[:no] if hash[:no]
    f.puts hash[:key] if hash[:key]
    f.puts hash[:material] if hash[:material]
    f.puts ((hash[:price] + 3) / 6.5).ceil(2)
    f.puts ((hash[:price] + 2.5) / 6.5).ceil(2)
    f.puts ((hash[:price] + 2) / 6.5).ceil(2)
    if hash[:note]
        ship_sku = hash[:note].split(",")
        ship_sku.each do |t|
            f.puts "第"+(t.to_i(10)-1).to_s+"个图片不要"
        end
     end
  end
end



    
