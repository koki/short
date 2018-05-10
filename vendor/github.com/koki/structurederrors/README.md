# structurederrors

Errors with context and other structured information

Errors in Golang lose context when one error leads to the calling function throwing an error, and its calling function throwing another error and so on.

Throwing errors using the functions provided in this library allow the user to capture the cascade of errors that are thrown. 

