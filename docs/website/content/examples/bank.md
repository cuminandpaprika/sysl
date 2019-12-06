---
title: "bank"
description: "A simple bank example"
date: 2019-11-08T14:37:15+11:00
weight: 20
draft: false
bref: ""
toc: true
---

``` yaml
# sysl ints -o "example2.svg" /example2 --project "Internet Banking" -v
# sysl ints -o "example2-epa.svg" /example2 --project "Integration" --epa -v
InternetBanking:
    GetCustomer:
        Customer <- GetCustomer
    GetAccountByCustomer:
        Ledger <- GetAccountByCustomer

Customer:
    GetCustomer: ...
Ledger:
    GetAccountByCustomer: ...

Integration [appfmt="%(appname)"]:
    _:
        InternetBanking
        Customer
        Ledger

```