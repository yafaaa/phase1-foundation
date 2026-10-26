==================
WARNING: DATA RACE
Read at 0x00c0000140f8 by goroutine 8:
  main.main.func1()
      /home/yafet/phase1-foundation/11-race-counter/main.go:18 +0x84

Previous write at 0x00c0000140f8 by goroutine 7:
  main.main.func1()
      /home/yafet/phase1-foundation/11-race-counter/main.go:18 +0x96

Goroutine 8 (running) created at:
  main.main()
      /home/yafet/phase1-foundation/11-race-counter/main.go:15 +0x7b

Goroutine 7 (finished) created at:
  main.main()
      /home/yafet/phase1-foundation/11-race-counter/main.go:15 +0x7b
==================
finished
Found 1 data race(s)
exit status 66

## Solution:
counter++ is not an atomic operation. 

    1. Reads the variable from RAM
    2. Arithmetic (Increment)
    3. update/write back to RAM

when goroutines are truely parallel, in this case, goroutine is at one of three steps(haven't finished writing back to RAM), while goroutine 8 is reading, which results data race/race condition.