# Go issue 22315 — outside evidence

Source: https://github.com/golang/go/issues/22315
Title: os: StartProcess ETXTBSY race on Unix systems
Author: rsc; opened 2017-10-18. Obtained through the web open tool on 2026-10-02.

The issue reports a Linux reproducer where concurrent threads write and execute scripts. A writable descriptor can survive in another thread's forked child until that child executes, despite close-on-exec. An earlier execution of the written script can consequently receive ETXTBSY. This is published outside evidence for the race mechanism; it does not itself prove this Spec's converted fixtures passed Linux CI.

The historical Actions log independently observes the Doctor failure; its adapter-lineage message does not directly establish ETXTBSY.
