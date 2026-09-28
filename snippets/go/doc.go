// Package snippets 是题解辅助函数的权威版本，只供复制，不供 import。
//
// 题解本体必须能原样提交到 LeetCode，而提交框只收一个文件，所以题解里用到的
// 辅助函数只能整段复制进去，不能引用这里。本目录的作用是让这些副本有一个
// 经过测试的源头，并由 check_test.go 保证副本不漂移：
//
//   - 这里每个非测试文件里的顶层函数（不含方法）都是一个 snippet，按函数名识别。
//   - leetcode/ 下任何题解定义了同名顶层函数，就视为引用了它，
//     从 func 关键字到函数结尾必须和这里逐字一致；函数上方的文档注释不比较，
//     各题可以写自己的说明。
//   - 这里只能 import 标准库，否则复制到 LeetCode 编译不过。
//   - 每个 snippet 都要在同名 _test.go 里有测试。
//
// 改 snippet 时，check_test.go 会列出所有不一致的副本，逐个同步即可。
// 包名是 snippets 而目录叫 go，和 structures/go 同理（go 是关键字）。
package snippets
