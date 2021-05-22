Within given four hours, I have complete below tasks

// Completed Task
```
* The program outputs a list of incorrectly spelled words. 
* The program includes the line and column number of the misspelled word
* The program handles proper nouns (person or place names, for example) correctly.
* For each misspelled word, the program outputs a list of suggested words.
```

// Not Completed Task
```
* The program prints the misspelled word along with some surrounding context.
```


// Setps to clone repo to your go src
```
- mkdir github
- cd github
- git clone https://github.com/ssenthil416/go-samples.git
- go get ./...
- cd go-samples/interview/wm/spellcheck
```

// Build and Run
- go run main.go


// Sample output
```
Line Number :7, Column Number :1, Wrong word :wecome, Suggested Word:
```


// Assumption
```
- Input file to validate is a text file.
- Nouns can be added to nouns file and if we run, the added nouns will be avoided from wrong list
```

Note: the words which are not part of dictionary are deplayed here too.
