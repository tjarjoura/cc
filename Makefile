build:
	go build ./cmd/cc
test: build
	./cc playground/test.c > test.asm
	nasm -f elf64 test.asm && gcc -no-pie test.o -o test
	- ./test
tccdump:
	tcc -c playground/test.c && objdump -Mintel -d test.o
tccrun:
	tcc -o test playground/test.c && ./test
	gcc -o test playground/test.c && ./test
gccrun:
clean:
	rm test.asm test.o test a.out
