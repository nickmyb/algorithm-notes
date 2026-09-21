// ===== 本地接线区 =====
// LeetCode 的 Java 编辑器预置了 java.util.*，本地要显式 import。
import java.util.ArrayList;
import java.util.List;

// ===== 以下是题解本体，与提交到 LeetCode 的代码一字不差 =====

class Solution {
    public List<Integer> inorderTraversal(TreeNode root) {
        return inorder(root, new ArrayList<>());
    }

    private List<Integer> inorder(TreeNode root, List<Integer> input) {
        if (root == null) {
            return input;
        }

        List<Integer> output = inorder(root.left, input);
        output.add(root.val);
        output = inorder(root.right, output);

        return output;
    }
}
