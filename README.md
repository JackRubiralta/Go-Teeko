TODO 
- redo encoder
- dont need offsets have it calcuate (have offset function like rankCombination but like rankOffset or rank calcutalteOffeset or unRankoffSet or unrank idk and improve code )
- rankCounts and unRankCounts (I think instead of offsets) and then rankPositions unRankPositions (instead of rankCombination) (THIS IS WHAT WE DO)
- For teeko improve it so the thigns could do any size game instead of just 5x5 but have it make the bit maps before
- Rename player_positions to player_bitmask
- 

To play 
```sh
go run main.go teeko.go encoder.go helper.go solver.go; ./main
```

To generate book
```sh
go build solver.go encoder.go teeko.go helper.go; ./solver
```

To unzip computed book
```sh
tar -xf book.zip
```

Now also do Teeko 78 

and also Teeko with 5 pieces 
1.1 Gb table 

Maybe improve book.txt to be 96 Mb
