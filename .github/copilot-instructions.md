# README Files

Do not include code in README files.
This requires more effort to keep the README's up to date.
Please instead provide a link to any code that should be referenced.

# Code Quality
- ALWAYS remove unused whitespace

# Code Structure

Please prefer moving types that have their own methods/behaviors to their own files.
The preference is towards a larger number of smaller files that fully encapsulate one concept.

# Go Code

Please make sure each go file that is created or edited if formated with `go fmt`.

Only get environment variables from the `main` package.  
Only exit from `main` package.

Make sure to `go mod tidy && go mod vendor` when packages are added or removed.

## Tests

Please prefer using the `stretchr/testify/suite` package for tests.

Do not create code just for the use of tests.  
If code only exists to make the tests work, then assume this is an error.  
If the existing code isn't easily testable, assume this is a problem with the code.

# UI Code

## PixiJS

When working with PixiJS code in the `ui/` directory, reference the official PixiJS documentation for LLMs:

- Full documentation: https://pixijs.com/llms-full.txt
- Medium documentation: https://pixijs.com/llms-medium.txt
- Documentation index: https://pixijs.com/llms.txt

These files are automatically updated daily by the PixiJS team and contain the complete API reference, examples, and best practices for PixiJS v8.

## UI Modules

### Specs

For each module please maintain a 'spec' folder that tracks the following details:
1. Technical Decisions
2. Features
3. Supporting documentation

The goal is for the files in this spec folder to guide future development, or any need to recreate module functionality.
