---
id: format-diagram
title: Diagram Format Args
---

import useBaseUrl from '@docusaurus/useBaseUrl';

## Format Arguments

By default, sequence diagrams generated by `sysl` only show the data type that is returned by an endpoint. You can instruct `sysl` to show the arguments to your endpoint in a sequence diagram.

Command:

```bash
sysl sd -o 'call-login-sequence.png' --epfmt '%(epname) %(args)' -s 'MobileApp <- Login' /assets/call.sysl -v call-login-sequence.png
```

See <a href={useBaseUrl('img/sysl/args.sysl')} >args.sysl</a> for complete example.

<img alt="Sequence Diagram" src={useBaseUrl('img/sysl/args-Seq.png')} />

A bit more explanation is required regarding `epname` and `args` keywords that are used in the `epfmt` command line argument.

## Display built-in attributes in appfmt and epfmt

`appfmt` and `epfmt` (app and endpoint format respectively) can be passed to the
`sysl sequence` and `sysl integrations` commands. They control how the application or endpoint name is
rendered as text.

There default value is `%(appname)` and `%(epname)` respectively.

These internal attributes are:

| Attribute  | Description                       |
| ---------- | --------------------------------- |
| appname    | short name of the application     |
| epname     | short name of the endpoint        |
| eplongname | Long quoted name of the endpoint  |
| controls   | controls defined on your endpoint |

### Example

```
App "Descriptive Long Application name":
  Endpoint-1 "Descriptive Long name for Endpoint 1":
    ...
  Endpoint-2 "Descriptive Long name for Endpoint 2":
    ...
```

Where:

| Attribute  | Value                                                                            |
| ---------- | -------------------------------------------------------------------------------- |
| appname    | App                                                                              |
| epname     | Endpoint-1 or Endpoint-2                                                         |
| eplongname | "Descriptive Long name for Endpoint 1" or "Descriptive Long name for Endpoint 2" |
| controls   | Controls defined on your endpoint                                                |

You can also refer to the attributes that you added by using `[]` or the
Collector syntax.

:::caution
TODO: Add more details about controls and an example
TODO: Add more details about the Collector syntax
:::

## Display custom attributes in fmt

You can display custom attributes in `epfmt` or `appfmt` arguments in the following
ways:

| Attribute                                                    | Value                                                       |
| ------------------------------------------------------------ | ----------------------------------------------------------- |
| %(@attrib_name)                                              | use `@` to refer to attrib_name                             |
| %(@attrib_name? yes_stmt \| no_stmt)                         | use the ternary operator `?` to test for existence of value |
| %(@attrib_name=='some_value'? yes_stmt &#124; no_stmt)       | compare attrib's value to some constant                     |
| %(@attrib_name=='a'? yes_stmt \|; @attrib_name=='b'? \| ...) | nested checks                                               |

The `stmt` can be any of the following types:

- plain-text: will be copied as-is to the output
- `<color red>text or %(attrib_name)</color>`: use html like syntax to color the
  output.

### Example

```html
appfmt="%(@status?<color red>%(appname)</color>|%(appname))" epfmt="%(@status?
<color green>%(epname)</color>|%(epname))"
```

See <a href={useBaseUrl('img/sysl/attribs.sysl')} >attribs.sysl</a> for complete example. Notice how
`appfmt` and `epfmt` use `%(@status)`.

<img alt="Sequence Diagram" src={useBaseUrl('img/diagramgen/attribs-Seq.png')} />