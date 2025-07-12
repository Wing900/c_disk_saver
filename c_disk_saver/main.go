package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

func main() {
	ensureAdmin()

	fmt.Println("--- C盘空间拯救向导 ---")
	fmt.Println("本工具将通过“移动文件夹并创建链接”的方式，为您的C盘瘦身。")
	fmt.Println("--------------------------------------------------")

	fmt.Println("\n【！！！极度危险警告 - 请在继续前仔细阅读！！！】")
	fmt.Println("1. 【绝对禁止】尝试移动核心系统文件夹，例如 'C:\\Windows' 或 'C:\\Program Files'！")
	fmt.Println("   这样做会导致您的系统立即崩溃且无法恢复！")
	fmt.Println("2. 【强制要求】在继续操作之前，您【必须】手动创建一个 Windows 系统还原点。")
	fmt.Println("   如果出现任何意外，这是您恢复系统的唯一希望。")
	fmt.Println("3. 【权限警告】即使以管理员身份运行，部分文件（特别是杀毒软件或系统底层文件）")
	fmt.Println("   也可能因为更深层的权限保护而导致迁移失败。")
	fmt.Println("\n【免责声明】本程序的作者不对因使用或不当使用本工具而造成的任何数据丢失、")
	fmt.Println("系统损坏或其他直接或间接损失承担任何责任。使用本工具即代表您同意自行承担所有风险。")
	fmt.Println("--------------------------------------------------")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("我已阅读并完全理解以上警告，并确认已创建系统还原点。同意继续吗？ (y/n): ")
	if !confirm(reader) {
		fmt.Println("操作已由用户中止。程序退出。")
		os.Exit(0)
	}

	fmt.Print("\n=> 请输入要移动的【源文件夹】的完整路径 (例如: C:\\Users\\您的用户名\\AppData\\Local\\Google):\n> ")
	sourcePath := getInput(reader)
	if !dirExists(sourcePath) {
		fmt.Printf("\n【错误】源文件夹未找到或路径错误: %s\n", sourcePath)
		pauseAndExit()
	}

	fmt.Print("\n=> 请输入【目标位置】的父文件夹路径 (例如: D:\\迁移的文件夹):\n> ")
	destParentPath := getInput(reader)
	if !dirExists(destParentPath) {
		fmt.Printf("\n【提示】目标位置的父文件夹不存在，正在尝试自动创建...\n")
		if err := os.MkdirAll(destParentPath, os.ModePerm); err != nil {
			fmt.Printf("【错误】创建目标文件夹失败: %v\n", err)
			fmt.Println("请检查磁盘权限或路径是否正确。")
			pauseAndExit()
		}
		fmt.Printf("【成功】已创建目标文件夹: %s\n", destParentPath)
	}

	destPath := filepath.Join(destParentPath, filepath.Base(sourcePath))

	fmt.Println("\n--- 请仔细核对以下信息 ---")
	fmt.Printf("【源文件夹】将被移动: %s\n", sourcePath)
	fmt.Printf("【目标文件夹】将存放在: %s\n", destPath)
	fmt.Println("--------------------------")
	fmt.Print("确认无误并开始执行吗？ (y/n): ")
	if !confirm(reader) {
		fmt.Println("操作已取消。")
		return
	}

	fmt.Println("\n--- 准备工作：请确保相关程序已关闭 ---")
	fmt.Println("为了防止文件被占用导致迁移失败，请手动关闭可能正在使用源文件夹的程序。")
	fmt.Printf("例如，如果您正在移动浏览器数据，请先【完全退出】该浏览器。\n")
	fmt.Println("关闭相关程序后，请按【回车键】继续...")
	bufio.NewReader(os.Stdin).ReadBytes('\n') // 等待用户按回车

	fmt.Println("\n--- 开始执行文件迁移 (第1步/共2步) ---")
	fmt.Println("正在移动文件... 这个过程可能会非常非常慢，请耐心等待，不要关闭窗口！")
	if !runCommand("robocopy", sourcePath, destPath, "/E", "/MOVE", "/XJ") {
		fmt.Println("\n【错误】文件迁移阶段失败。请检查上方的 Robocopy 输出日志，查找具体的错误原因。")
		fmt.Println("失败的常见原因包括：文件被其他程序占用、目标磁盘空间不足等。")
		pauseAndExit()
	}
	fmt.Println("\n--- 文件迁移阶段已完成。 ---")

	fmt.Println("\n--- 开始创建目录链接 (第2步/共2步) ---")
	for {
		if runCommand("cmd", "/C", "mklink", "/J", sourcePath, destPath) {
			fmt.Println("\n【🎉 巨大成功!】操作已顺利完成！")
			fmt.Println("链接已创建，C盘空间已被释放。请现在去验证对应的软件是否能正常工作。")
			break // 成功后跳出循环
		}

		fmt.Println("\n【错误】创建链接失败！这通常意味着C盘的原始文件夹“空壳”没有被自动删除。")
		fmt.Println("=> 现在，请您【手动删除】那个空的源文件夹:")
		fmt.Printf("   路径: %s\n", sourcePath)
		fmt.Print("   删除后，请按 'y' 重试，或按 'n' 退出程序。 (y/n): ")

		if !confirm(reader) {
			fmt.Println("操作已中止。如果文件已移动，您可能需要手动创建链接。")
			break // 用户选择退出
		}
	}

	pauseAndExit()
}

// runCommand 执行一个外部命令，并将其标准输出和标准错误实时显示到控制台。
// 特别处理了 Robocopy 的退出码逻辑。
func runCommand(name string, args ...string) bool {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err == nil {
		return true // 退出码为 0，总是表示成功
	}

	// 检查错误类型是否为 ExitError，这表示命令已执行但退出码非 0
	if exitError, ok := err.(*exec.ExitError); ok {
		// 获取具体的退出码
		exitCode := exitError.ExitCode()

		// 针对 robocopy 的智能判断：小于8的退出码都表示成功或部分成功
		if name == "robocopy" && exitCode < 8 {
			fmt.Printf("\n【提示】Robocopy 操作已完成，退出码为 %d (这表示文件已成功迁移)。\n", exitCode)
			return true
		}
	}

	// 对于所有其他错误（包括 robocopy >= 8 的情况或命令无法执行的错误）
	fmt.Printf("\n【底层错误】执行命令 '%s' 失败: %v\n", name, err)
	return false
}

func getInput(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func confirm(reader *bufio.Reader) bool {
	input := getInput(reader)
	return strings.ToLower(input) == "y"
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func pauseAndExit() {
	fmt.Println("\n按回车键退出程序...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(1)
}

// ensureAdmin 检查管理员权限，如果不足则尝试提权并重新运行
func ensureAdmin() {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		runAsAdmin()
		os.Exit(0)
	}
}

// runAsAdmin 使用 Windows API 的 'runas' 动作来请求UAC提权
func runAsAdmin() {
	verb := "runas"
	exe, _ := os.Executable()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	argsPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 // SW_SHOWNORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argsPtr, nil, showCmd)
	if err != nil {
		fmt.Printf("请求管理员权限失败: %v\n", err)
	}
}
