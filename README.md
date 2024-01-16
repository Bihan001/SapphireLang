## Backus-Naur Form (BNF) grammar

```
<program> ::= <statement> | <statement> <program>

<statement> ::= <function_definition> | <variable_declaration> | <expression> SEMICOLON

<function_definition> ::= INT <identifier> LPAREN <parameters> RPAREN LBRACE <program> RBRACE

<parameters> ::= INT <identifier> | INT <identifier> COMMA <parameters> | ε

<variable_declaration> ::= INT <identifier> SEMICOLON | INT <identifier> EQUAL <expression> SEMICOLON

<expression> ::= <term> | <expression> PLUS <term> | <expression> MINUS <term>

<term> ::= <factor> | <term> STAR <factor> | <term> SLASH <factor>

<factor> ::= <identifier> | <integer> | LPAREN <expression> RPAREN | <unary> | <ternary>

<unary> ::= MINUS <factor> | NOT <factor>

<ternary> ::= <expression> QUESTION <expression> COLON <expression>

<identifier> ::= ID

<integer> ::= NUM
```

