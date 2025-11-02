# README Files

Do not include code in README files.
This requires more effort to keep the README's up to date.
Please instead provide a link to any code that should be referenced.

# Code Structure

Please prefer moving types that have their own methods/behaviors to their own files.
The preference is towards a larger number of smaller files that fully encapsulate one concept.

# Go Code

Please make sure each go file that is created or edited if formated with `go fmt`.

Only get environment variables from the `main` package.  

## Tests

Please prefer using the `stretchr/testify/suite` package for tests.

Do not create code just for the use of tests.  
If code only exists to make the tests work, then assume this is an error.  
If the existing code isn't easily testable, assume this is a problem with the code.  
