build:
	go build ./cmd/cc
test: build
	./cc playground/test.c > test.asm
	nasm -f elf64 test.asm && gcc -no-pie test.o -o test
	- ./test
clean:
	rm test.asm test.o test a.out
