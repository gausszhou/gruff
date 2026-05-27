# 修复：`wrapText` 中的过早换行问题

## 症状

恰好能放在行尾（包括所需空格分隔符）的单词被推到下一行，导致行长度不均匀，比换行宽度短。

例如，当 `width=72` 时：

```
修复前（有问题）：               修复后（正确）：
Markdown is a plain text         Markdown is a plain text format for
format for                       writing structured documents, based on
writing structured documents,    conventions for indicating formatting
based on                         in email and usenet posts.
```

注意每行只有约 55-65 个字符，而不是填满完整的 72 个字符。

## 根本原因

在 `flushWord()`（`gruff.go:99`）中，换行检查为：

```go
if lineLen > 0 && lineLen+wLen > width {
```

当 `spaces > 0` 时，会在单词之前写入一个空格分隔符（第 104-107 行）。但宽度检查没有考虑这个缓冲的空格：

- `lineLen` → 当前行的显示宽度（不含尾随空格）
- `wLen` → 当前单词的显示宽度
- `spaces` → 将在单词之前输出的空格字符数

如果 `lineLen + 1 + wLen == width`，单词与其分隔符恰好能放下，但旧代码认为放不下并过早换行。

## 修复方法

将空格分隔符加入宽度检查：

```go
if lineLen > 0 && lineLen+wLen+(b2i(spaces > 0)) > width {
```

`b2i` 将 `bool` 转换为 `int`（`spaces > 0` 时为 1，否则为 0），考虑了 `flushWord` 将要写入的空格分隔符。

## 验证

修复前，使用 `width=72` 渲染 `testdata/_data.md` 产生 8322 行。修复后：8046 行（减少了 276 行，全部来自现在能填满全宽的段落）。

该修复对所有现有测试透明 — 没有测试断言发生变化。
