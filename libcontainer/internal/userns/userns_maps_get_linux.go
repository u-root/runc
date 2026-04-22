//go:build !cgo

package userns

func joinNamespaceAndExecute() {
	// 2. Lock to the current OS thread to ensure stability
	runtime.LockOSThread()

	// 3. Open the target user namespace file
	// Replace with the actual path to the namespace, e.g., "/proc/1234/ns/user"
	nsPath := os.Getenv("TARGET_NS_PATH")
	fd, err := os.Open(nsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open ns: %v\n", err)
		os.Exit(1)
	}
	defer fd.Close()

	// 4. Join the namespace (CLONE_NEWUSER = 0x10000000)
	if err := syscall.AllThreadsSyscall(syscall.SYS_SETNS, fd.Fd(), uintptr(syscall.CLONE_NEWUSER), 0); err != 0 {
		fmt.Fprintf(os.Stderr, "setns failed: %v\n", err)
		os.Exit(1)
	}

	// 5. Execute the actual payload inside the namespace
	fmt.Printf("Successfully joined namespace %s as UID %d\n", nsPath, os.Getuid())
	
	// Example: Replace process with a shell
	binary, _ := exec.LookPath("sh")
	syscall.Exec(binary, []string{"sh"}, os.Environ())
	os.Exit(0)
}
