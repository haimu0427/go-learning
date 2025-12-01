--学习lua是为了在nginx里使用lua脚本
-- Lua是一种轻量级的脚本语言，设计目标是嵌入式应用程序中使用。
print("Hello, World!")
-- 以下划线开头连接一串大写字母的名字（比如 _VERSION）
--[[变量总是认为是全局的。全局变量不需要声明，
给一个变量赋值后即创建了这个全局变量，访问一个没有初始化的全局变量也不会出错，只不过得到的结果是：nil。
如果你想删除一个全局变量，只需要将变量赋值为nil。
当且仅当一个变量不等于nil时，这个变量即存在。]]--

local str = "hello lua"
local num = 29
local bool = true
-- table类型
local arr = {'java', 'c++', 'lua'}
local map = {name = "tom", age = 20}
print(arr[1]) -- 访问数组, 注意lua的数组下标从1开始
print(map.name) -- 访问map
local str_concat = str .. "!" -- 字符串连接符
print(str_concat)
for index, value in ipairs(arr) do 
    print(index, value)
end
for key, value in pairs(map) do
    print(key, value)
end
--在对一个数字字符串上进行算术操作时，Lua 会尝试将这个数字字符串转成一个数字:
-- 和java有很大的不同, java中偏向与向着string类型转换, lua偏向与向着number类型转换

--  # 来计算字符串的长度，放在字符串前面，
print(#str)
function factorial(n)
    if n == 0 then 
        return 1
    else
        return n * factorial(n - 1)
    end
end
print(factorial(5))
factorial2 = factorial
print(factorial2(6))
-- 多返回值
