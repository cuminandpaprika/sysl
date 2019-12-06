---
title: "Generating Go With Sysl"
description: "A walkthrough for generating Go from Sysl"
date: 2019-11-08T11:36:37+11:00
draft: false
toc: false
---

## Sysl Go Code Generation

Sysl codegen works by taking in a model which describes the data types and endpoints of your application.

It interprets the model file which is written in SYSL using the grammar file
It then reads the transform files and outputs a go file for each, using a golang grammar file to translate the output.
Transforms are written in a sysl-transformation language.

For example, one output file might be a client library, which can be used to hit the endpoints exposed by the application. 
Another output might be the server library that exposes the application's API endpoints. 

## Sysl Model

Let's start with the following sysl model. It defines a simple Todo web application.

``` yaml

Todos:
  !type Todo:
    userId <: int
    id <: int
    title <: string
    completed <: bool


  !type Post:
    userId <: int
    id <: int
    title <: string
    body <: string

  !alias Posts:
    sequence of Post

  !type ErrorResponse:
    status <: string

  !type ResourceNotFoundError:
    status <: string

  /todos:
    /{id<:int}:
      GET:
        if notfound:
          return 404 <: ResourceNotFoundError
        else if failed:
          return 500 <: ErrorResponse
        else:    
          return 200 <: Todo

  /posts:
    GET:
      if notfound:
        return 404 <: ResourceNotFoundError
      else if failed:
        return 500 <: ErrorResponse
      else:    
        return 200 <: Posts

  /comments:
    GET ?postId=int:
      return Posts
      
    POST (newPost <: Post [~body]):
      return Post
```

## Grammar

Now lets look at a grammar file.

Let's break down the first few lines of the grammar file:

* goFile has a PackageClause, Comment(optional), ImportDecl(optional) and a list of TopLevelDecl
* PackageClause contains PackageName
* ImportDecl has a list of ImportSpec
* ImportSpec can be just an Import definition or a NamedImport
* NamedImport has a Name and Import definition
* TopLevelDecl can be a Comment followed by either a Declaration, or a FunctionDecl or a MethodDecl

``` yaml
goFile: PackageClause Comment? '\n' ImportDecl? '\n' TopLevelDecl+ '\n';
PackageClause: 'package' PackageName '\n';

ImportDecl: 'import' '(\n' ImportSpec* '\n)\n';
ImportSpec: (Import | NamedImport) '\n';
NamedImport: Name Import;
TopLevelDecl: Comment '\n' (Declaration | FunctionDecl | MethodDecl);
Declaration: VarDecl | VarDeclWithVal | ConstDecl | StructType | InterfaceType | AliasDecl;
StructType : 'type' StructName 'struct' '{\n' FieldDecl* '}\n\n';
FieldDecl: '\t' identifier Type? Tag? '\n';
IdentifierList: identifier IdentifierListC*;
IdentifierListC: ',' identifier;

VarDeclWithVal: 'var' identifier '=' TypeName '\n';
VarDecl: 'var' identifier TypeName '\n';
ConstDecl: 'const' '(\n'  ConstSpec '\n)\n';
ConstSpec: VarName TypeName '=' ConstValue '\n';

FunctionDecl   : 'func' FunctionName Signature? Block '\n\n';
Signature: Parameters Result?;
Parameters: '(' ParameterList? ')';
Result         : ReturnTypes | TypeName;
ReturnTypes: '(' TypeName ResultTypeList* ')';
ResultTypeList: ',' TypeName ;
TypeList:  TypeName;
ParameterList     : ParameterDecl ParameterDeclC*;
ParameterDecl  : Identifier TypeName;
ParameterDeclC: ',' ParameterDecl;

InterfaceType      : 'type' InterfaceName 'interface'  '{\n'  MethodSpec* '}\n\n' MethodDecl*;
MethodSpec         : '\t' MethodName Signature '\n' | InterfaceTypeName ;
MethodDecl: 'func' Receiver FunctionName Signature? Block? '\n\n';
Receiver: '(' ReceiverType ')';
AliasDecl: 'type' identifier Type? ';\n\n';

Block: '{\n'  StatementList* '}\n';
StatementList: '\t' Statement '\n';
Statement: ReturnStmt |  DeclareAndAssignStmt | AssignStmt | IfElseStmt | IncrementVarByStmt | FunctionCall | VarDecl;

AssignStmt: Variables '=' Expression;
IfElseStmt: 'if' Expression Block;
IncrementVarByStmt: Variables '+=' Expression;
ReturnStmt: 'return' (PayLoad | Expression);
DeclareAndAssignStmt: Variables ':=' Expression;

Expression: FunctionCall | NewStruct | GetArg |  ValueExpr | NewSlice;

GetArg: LHS '.' RHS;
NewSlice: '[]' TypeName '{' SliceValues? '}';
FunctionCall: FunctionName '(' FunctionArgs? ')';
FunctionArgs: Expression FuncArgsRest*;
FuncArgsRest: ',' Expression;
NewStruct: StructName '{}';
```

## Transforms

Now let's look at a simple transform. At first glance, it might look similar to the SYSL model.

Lets break this down:

* In the first line, `CodeGenTransform:` defines the name of the transform.
* `!view filename(app <: sysl.App) -> string:` defines a view which takes in an input of type `sysl.App` and returns a `string`
* `!view goFile(app <: sysl.App) -> string:` this is where the the contents of the output Go file are defined

A transform needs 2 things to be correctly interpreted by sysl

* A view called 'filename' must exist in the transform as it defines the name of output go file.
* The entry point should be defined as a sysl view with sysl.App as the argument. Also, the entry point should be provided as the -start argument in syslgen command.

``` yaml
CodeGenTransform:
  !view filename(app <: sysl.App) -> string:
    app -> (:
      filename =  "types.go"
    )

  !view goFile(app <: sysl.App) -> string:
    app -> (:
      # Rest of transform
    )
```

Navigate to syslgen-examples directory using command line and use following syslgen command is used to generate code.

syslgen -grammar grammars/go.gen.g -model examples/todos.sysl -outdir . -root-model . -root-transform . -start goFile -transform transforms/todo-transform.sysl
At this stage, it is obvious that the transformation is invalid as mandatory keywords are still missing in the transform. Therefore, above command will give an erroneous output.



In following example, PackageClause and its subsequent keywords are defined. Note that f does not have subsequent keywords. That means, it can be assigned with free text that defines the package name. Same with the Comment. For comment, two forward slashes has been added so that it will be a comment in generated go file.

``` yaml
CodeGenTransform:
  !view filename(app <: sysl.App) -> string:
    app -> (:
      filename =  "types.go"
    )

  !view goFile(app <: sysl.App) -> string:
    app -> (:
      PackageClause = app -> <PackageClause>(:
        PackageName = "todo"
      )


      Comment = "// This is a comment"

      TopLevelDecl = app -> <TopLevelDecl>(:
        # Top level declarations
      )
    )
```

Before diving into top level declarations in todos.sysl file, let's make a simple definition of it first so that sysl transformation compile successfully and gives an output. In that way, it will be easier to understand syntaxes and concepts in writing transformation.
```yaml
CodeGenTransform:
  !view filename(app <: sysl.App) -> string:
    app -> (:
      filename =  "types.go"
    )


  !view goFile(app <: sysl.App) -> string:
    app -> (:
      PackageClause = app -> <PackageClause>(:
        PackageName = "todo"
      )

      Comment = "//\n//    THIS IS AUTOGENERATED BY syslgen \n//\n"

      TopLevelDecl = [app] -> <TopLevelDecl>(:
        Comment = "// Struct"
        Declaration = app -> <Declaration>(:
          StructType = app -> <StructType>(:
            StructName = "a"
            FieldDecl = [app] -> <FieldDecl>(:
              identifier = "b"
              Type = "c"
            )
          )
        )
      )
    )
```

Notice that in the first line of grammar, it says "goFile has a PackageClause, Comment(optional), ImportDecl(optional) and a list of TopLevelDecl". Therefore, in transformation, TopLevelDecl should be defined as a list. In order to achieve that, 'app' is surrounded with square brackets. Same done for 'FieldDecl' as well because StructType expects a list of FieldDecl. In sysl, if a variable is surrounded with square brackets, it begins to behave as a list. As the result, sysl begins to iterate over the list which has only one element in this example and returns a list with one element as the result. This method is the simplest way of obtaining a list. Alternatively, it is possible to get single value instead of a list, assign it to a variable, and convert it to a list later. See following example.

```yaml
CodeGenTransform:
  !view filename(app <: sysl.App) -> string:
    app -> (:
      filename =  "types.go"
    )


  !view goFile(app <: sysl.App) -> string:
    app -> (:
      PackageClause = app -> <PackageClause>(:
        PackageName = "todo"
      )

      Comment = "//\n//    THIS IS AUTOGENERATED BY syslgen \n//\n"

      let singleTopLevelDecl = app -> <TopLevelDecl>(:
        Comment = "// Struct"
        Declaration = app -> <Declaration>(:
          StructType = app -> <StructType>(:
            StructName = "a"
            FieldDecl = [app] -> <FieldDecl>(:
              identifier = "b"
              Type = "c"
            )
          )
        )
      )

      TopLevelDecl = [singleTopLevelDecl]
    )

```

There are two type definitions in todos.sysl file as 'Todo' and 'Post'. The target is to create two struct declarations for those types. All type definitions can be retrieved as a set using 'app.types'. Therefore, a list of TopLevelDecl can be created iterating app.types. As the first step, replace '[app]' with 'app.types'. As 'app.types' is already a type of collection, it is not required to surround it with square brackets.

Now sysl transformation is iterating over types and creating a list of TopLevelDecl. This will give two identical struct types as StructName, identifier, and Type are still hard coded.

When iterating over a collection, it is required to define the iterator variable. The following example, 'type' is used as the variable name.

TopLevelDecl = app.types -> <TopLevelDecl>(type:

```yaml
CodeGenTransform:
  !view filename(app <: sysl.App) -> string:
    app -> (:
      filename =  "types.go"
    )


  !view goFile(app <: sysl.App) -> string:
    app -> (:
      PackageClause = app -> <PackageClause>(:
        PackageName = "todo"
      )

      Comment = "//\n//    THIS IS AUTOGENERATED BY syslgen \n//\n"

      TopLevelDecl = app.types -> <TopLevelDecl>(type:
        Comment = "// Struct"
        Declaration = app -> <Declaration>(:
          StructType = app -> <StructType>(:
            StructName = "a"
            FieldDecl = [app] -> <FieldDecl>(:
              identifier = "b"
              Type = "c"
            )
          )
        )
      )
    )
```

At this point, name of the structs and field definitions of each struct should be read in order to generate proper struct definitions for types defined in todos.sysl.

Name of the struct can be obtained using 'type.key'. Collection of Field declarations are stored in 'type.value.fields'. Type of each field should be matched to corresponding golang type. Following is the complete transform for generating structs in golang using sysl file.

```yaml
CodeGenTransform:
  !view filename(app <: sysl.App) -> string:
    app -> (:
      filename =  "types.go"
    )

  !view goFile(app <: sysl.App) -> string:
    app -> (:
      PackageClause = app -> <PackageClause>(:
        PackageName = "todo"
      )

      Comment = "//\n//    THIS IS AUTOGENERATED BY syslgen \n//\n"

      TopLevelDecl = app.types -> <TopLevelDecl>(type:
        Comment = "// Struct"
        Declaration = type -> <Declaration>(:
          StructType = type -> <StructType>(:
            StructName = type.key
            FieldDecl = type.value.fields -> <FieldDecl>(field:
              identifier = field.key
              Type = if field.value.type ==:
                "primitive" => if field.value.primitive ==:
                  "DECIMAL" => "double"
                  "INT" => "int64"
                  "FLOAT" => "float64"
                  "STRING" => "string"
                  "BOOL" => "bool"
            )
          )
        )
      )
    )
```

When moving forward, it will be required to convert sysl types to relevant golang types in multiple places. Therefore, it is a good practice to define that transformation as a separate sysl view so that it can be called as a function whenever type conversion requires. Further, another view can be defined to obtain all struct definitions using a single call so that goFile view will be more readable and will not get bulky.

```yaml
CodeGenTransform:
  !view filename(app <: sysl.App) -> string:
    app -> (:
      filename =  "types.go"
    )

  !view GoType(t <: sysl.Type) -> string:
    t -> (:
      let IsPtr = if t.optional == true && t.type != "sequence" then "*" else ""
      let typeName = if t.type ==:
        "primitive" => if t.primitive ==:
          "DECIMAL" => "double"
          "INT" => "int64"
          "FLOAT" => "float64"
          "STRING" => "string"
          "STRING_8" => "string"
          "BOOL" => "bool"
          "DATE" => "date.Date"
          "DATETIME" => "time.Time"
        "sequence" => "[]" + GoType(t.sequence).out
        else if HasPrefix(t.type_ref, "EXTERNAL_") == true then GoName(Trim(t.type_ref, "EXTERNAL_")).out else GoName(t.type_ref).out
      out = IsPtr + typeName
    )
  
  !view getStructDefs(types <: set of Type) -> sequence of TopLevelDecl:
    types -> (type:
        let typeName = type.key
        Comment = '// ' + typeName + ' ...'
        Declaration = type -> <Declaration>(:
          StructType = type -> <StructType>(:
            StructName = typeName
            FieldDecl = type.value.fields -> <FieldDecl>(field:
              identifier = field.key
              Type = GoType(field.value).out
            )
          )
        )
    )

  !view goFile(app <: sysl.App) -> string:
    app -> (:
      PackageClause = app -> <PackageClause>(:
        PackageName = "todo"
      )

      Comment = "//\n//    THIS IS AUTOGENERATED BY syslgen \n//\n"

      TopLevelDecl = getStructDefs(app.types)
    )
```
