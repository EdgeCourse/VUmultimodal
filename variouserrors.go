/*
Error handling in Go is a little different than other mainstream programming languages like Java, JavaScript, or Python. Go's built-in errors don't contain stack traces, nor do they support conventional try/catch methods to handle them. Instead, errors in Go are just values returned by functions, and they can be treated in much the same way as any other datatype - leading to a surprisingly lightweight and simple design.
*/


//The Error Type

//The error type in Go is implemented as the following interface:

/*
type error interface {
    Error() string
}
*/

//an error is anything that implements the Error() method, which returns an error message as a string. 

//Constructing Errors

/*
Errors can be constructed on the fly using Go's built-in errors or fmt packages. For example, the following function uses the errors package to return a new error with a static error message:
*/

/*
import "errors"

func DoSomething() error {
    return errors.New("something didn't work")
}
*/

//Similarly, the fmt package can be used to add dynamic data to the error, such as an int, string, or another error. For example:



import "fmt"

func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("can't divide '%d' by zero", a)
    }
    return a / b, nil
}

//Note that fmt.Errorf will prove extremely useful when used to wrap another error with the %w format verb 


/*
Errors can be returned as nil, and in fact, it's the default, or "zero", value of on error in Go. This is important since checking if err != nil is the idiomatic way to determine if an error was encountered (replacing the try/catch statements you may be familiar with in other programming languages).

Errors are typically returned as the last argument in a function. Hence in our example above, we return an int and an error, in that order.

When we do return an error, the other arguments returned by the function are typically returned as their default "zero" value. A user of a function may expect that if a non-nil error is returned, then the other arguments returned are not relevant.

Lastly, error messages are usually written in lower-case and don't end in punctuation. Exceptions can be made though, for example when including a proper noun, a function name that begins with a capital letter, etc.
*/

/*
Defining Expected Errors

Another important technique in Go is defining expected Errors so they can be checked for explicitly in other parts of the code. This becomes useful when you need to execute a different branch of code if a certain kind of error is encountered.
*/

/*
Defining Sentinel Errors

Building on the Divide function from earlier, we can improve the error signaling by pre-defining a "Sentinel" error. Calling functions can explicitly check for this error using errors.Is:
*/


import (
    "errors"
    "fmt"
)

var ErrDivideByZero = errors.New("divide by zero")

func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, ErrDivideByZero
    }
    return a / b, nil
}

func main() {
    a, b := 10, 0
    result, err := Divide(a, b)
    if err != nil {
        switch {
        case errors.Is(err, ErrDivideByZero):
            fmt.Println("divide by zero error")
        default:
            fmt.Printf("unexpected division error: %s\n", err)
        }
        return
    }

    fmt.Printf("%d / %d = %d\n", a, b, result)
}

/*
Defining Custom Error Types

Many error-handling use cases can be covered using the strategy above, however, there can be times when you might want a little more functionality. Perhaps you want an error to carry additional data fields, or maybe the error's message should populate itself with dynamic values when it's printed.

You can do that in Go by implementing custom errors type.

Below is a slight rework of the previous example. Notice the new type DivisionError, which implements the Error interface. We can make use of errors.As to check and convert from a standard error to our more specific DivisionError.
*/



import (
    "errors"
    "fmt"
)

type DivisionError struct {
    IntA int
    IntB int
    Msg  string
}

func (e *DivisionError) Error() string { 
    return e.Msg
}

func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, &DivisionError{
            Msg: fmt.Sprintf("cannot divide '%d' by zero", a),
            IntA: a, IntB: b,
        }
    }
    return a / b, nil
}

func main() {
    a, b := 10, 0
    result, err := Divide(a, b)
    if err != nil {
        var divErr *DivisionError
        switch {
        case errors.As(err, &divErr):
            fmt.Printf("%d / %d is not mathematically valid: %s\n",
              divErr.IntA, divErr.IntB, divErr.Error())
        default:
            fmt.Printf("unexpected division error: %s\n", err)
        }
        return
    }

    fmt.Printf("%d / %d = %d\n", a, b, result)
}

/*
Note: when necessary, you can also customize the behavior of the errors.Is and errors.As. See this Go.dev blog for an example.

Another note: errors.Is was added in Go 1.13 and is preferable over checking err == .... More on that below.

Wrapping Errors

In these examples so far, the errors have been created, returned, and handled with a single function call. In other words, the stack of functions involved in "bubbling" up the error is only a single level deep.

Often in real-world programs, there can be many more functions involved - from the function where the error is produced, to where it is eventually handled, and any number of additional functions in-between.
*/
/*
In Go 1.13, several new error APIs were introduced, including errors.Wrap and errors.Unwrap, which are useful in applying additional context to an error as it "bubbles up", as well as checking for particular error types, regardless of how many times the error has been wrapped.

A bit of history: Before Go 1.13 was released in 2019, the standard library didn't contain many APIs for working with errors - it was basically just errors.New and fmt.Errorf. As such, you may encounter legacy Go programs in the wild that do not implement some of the newer error APIs. Many legacy programs also used 3rd-party error libraries such as pkg/errors. Eventually, a formal proposal was documented in 2018, which suggested many of the features we see today in Go 1.13+.
*/
/*
The Old Way (Before Go 1.13)

It's easy to see just how useful the new error APIs are in Go 1.13+ by looking at some examples where the old API was limiting.

Let's consider a simple program that manages a database of users. In this program, we'll have a few functions involved in the lifecycle of a database error.

For simplicity's sake, let's replace what would be a real database with an entirely "fake=" database that we import from "example.com/fake/users/db".

Let's also assume that this fake database already contains some functions for finding and updating user records. And that the user records are defined to be a struct that looks something like:

fake/users/db.gopackage db
*/

/*
type User struct {
  ID       string
  Username string
  Age      int
}

func FindUser(username string) (*User, error) { /* ... */ }
func SetUserAge(user *User, age int) error { /* ... */ }
Here's our example program:


import (
    "errors"
    "fmt"

    "example.com/fake/users/db"
)

func FindUser(username string) (*db.User, error) {
    return db.Find(username)
}

func SetUserAge(u *db.User, age int) error {
    return db.SetAge(u, age)
}

func FindAndSetUserAge(username string, age int) error {
  var user *User
  var err error

  user, err = FindUser(username)
  if err != nil {
      return err
  }

  if err = SetUserAge(user, age); err != nil {
      return err
  }

  return nil
}

func main() {
    if err := FindAndSetUserAge("bob@example.com", 21); err != nil {
        fmt.Println("failed finding or updating user: %s", err)
        return
    }

    fmt.Println("successfully updated user's age")
}

*/

/*
Now, what happens if one of our database operations fails with some malformed request error?

The error check in the main function should catch that and print something like this:

>_failed finding or updating user: malformed request
But which of the two database operations produced the error? Unfortunately, we don't have enough information in our error log to know if it came from FindUser or SetUserAge.

Go 1.13 adds a simple way to add that information.
*/
Errors Are Better Wrapped

The snippet below is refactored so that is uses fmt.Errorf with a %w verb to “wrap” errors as they “bubble up” through the other function calls. This adds the context needed so that it's possible to deduce which of those database operations failed in the previous example.

main.gopackage main

import (
    "errors"
    "fmt"

    "example.com/fake/users/db"
)

func FindUser(username string) (*db.User, error) {
    u, err := db.Find(username)
    if err != nil {
        return nil, fmt.Errorf("FindUser: failed executing db query: %w", err)
    }
    return u, nil
}

func SetUserAge(u *db.User, age int) error {
    if err := db.SetAge(u, age); err != nil {
      return fmt.Errorf("SetUserAge: failed executing db update: %w", err)
    }
}

func FindAndSetUserAge(username string, age int) error {
  var user *User
  var err error

  user, err = FindUser(username)
  if err != nil {
      return fmt.Errorf("FindAndSetUserAge: %w", err)
  }

  if err = SetUserAge(user, age); err != nil {
      return fmt.Errorf("FindAndSetUserAge: %w", err)
  }

  return nil
}

func main() {
    if err := FindAndSetUserAge("bob@example.com", 21); err != nil {
        fmt.Println("failed finding or updating user: %s", err)
        return
    }

    fmt.Println("successfully updated user's age")
}
If we re-run the program and encounter the same error, the log should print the following:

>_failed finding or updating user: FindAndSetUserAge: SetUserAge: failed executing db update: malformed request
Now our message contains enough information that we can see the problem originated in the db.SetUserAge function. Phew! That definitely saved us some time debugging!

If used correctly, error wrapping can provide additional context about the lineage of an error, in ways similar to a traditional stack-trace.

Wrapping also preserves the original error, which means errors.Is and errors.As continue to work, regardless of how many times an error has been wrapped. We can also call errors.Unwrap to return the previous error in the chain.

Curious to learn how error wrapping works under the hood? Take a peek at the internals of fmt.Errorf, the %w verb, and the errors API.

When To Wrap

Generally, it's a good idea to wrap an error with at least the function's name, every time you “bubble it up” - i.e. every time you receive the error from a function and want to continue returning it back up the function chain.


Wrapping an error adds the gift of context
There are some exceptions to the rule, however, where wrapping an error may not be appropriate.

Since wrapping the error always preserves the original error messages, sometimes exposing those underlying issues might be a security, privacy, or even UX concern. In such situations, it could be worth handling the error and returning a new one, rather than wrapping it. This could be the case if you're writing an open-source library or a REST API where you don't want the underlying error message to be returned to the 3rd-party user.

While you're here:

Earthly is the effortless CI/CD framework. 
Develop CI/CD pipelines locally and run them anywhere!

Conclusion

That's a wrap! In summary, here's the gist of what was covered here:

Errors in Go are just lightweight pieces of data that implement the Error interface
Predefined errors will improve signaling, allowing us to check which error occurred
Wrap errors to add enough context to trace through function calls (similar to a stack trace)
I hope you found this guide to effective error handling useful. If you'd like to learn more, I've attached some related articles I found interesting during my own journey to robust error handling in Go.

References

Error handling and Go
Go 1.13 Errors
Go Error Doc
Go By Example: Errors
Go By Example: Panic
  <div data-element="header">
    <h4 class="text-2xl font-bold pb-2">Get notified about new articles!</h4>
  </div>
  <!-- <div class="formkit-subheader" data-element="subheader" style="color: rgb(104, 104, 104); font-size: 18px;">
Subscribe to get our latest content by email.

–>
    <div class="formkit-field">
      <input
        class="border rounded-md p-2 mb-2 w-2/3"
        name="email_address"
        aria-label="Email Address"
        placeholder="Email Address"
        required=""
        type="email"
      />
    </div>
    <button
      data-element="submit"
      class="formkit-submit formkit-submit w-44 pt-1 pb-1"
      style="
        color: rgb(255, 255, 255);
        background-color: rgb(22, 119, 190);
        border-radius: 4px;
        font-weight: 400;
      "
    >
      <div class="formkit-spinner">
        <div></div>
        <div></div>
        <div></div>
      </div>
      <span class="">Subscribe</span>
    </button>
  </div>
  <div class="formkit-guarantee" data-element="guarantee">
    <p class="text-slate-400 text-sm">
      We won't send you spam. Unsubscribe at any time.
    </p>
  </div>
</div>
<style></style>
