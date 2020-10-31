---
title: ":house: Home"
date: 2020-03-08
---

[![GitHub Action Build](https://github.com/tamada/lioss/workflows/build/badge.svg?branch=master)](https://github.com/tamada/lioss/actions?workflow=build)
[![Coverage Status](https://coveralls.io/repos/github/tamada/lioss/badge.svg?branch=master)](https://coveralls.io/github/tamada/lioss?branch=master)
[![codebeat badge](https://codebeat.co/badges/dc3481f5-852b-4537-a5f5-150e2bfa998c)](https://codebeat.co/projects/github-com-tamada-lioss-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/tamada/lioss)](https://goreportcard.com/report/github.com/tamada/lioss)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/tamada/lioss/blob/master/LICENSE)
[![Version](https://img.shields.io/badge/Version-1.0.0-yellowgreen.svg)](https://github.com/tamada/lioss/releases/tag/v1.0.0)

## :speaking_head: Overview

Generally, OSS projects have licenses.
The licenses grant permissions to users for using, modifying, and sharing the software.
The users of the software must follow the terms shown in the licenses.

On the other hand, today's software generally has some dependencies.
Additionally, dependant software has some dependencies, too.
Therefore, the dependant graph of the OSS becomes complex.

In such a situation, it is a quite tough task for checking the conflicts among licenses.
The first problem is to detect a conflict between two given licenses.
The second problem is to identify the license of a project.
`lioss` tries to solve the above second problem by identifying the license of the given project.

SPDX is trying to automatically identify licenses, however,  it is hard to say that it became common sense.
This project detects the OSS licenses from the LICENSE files of the given projects.
Then, we aim to detect conflicts by identifying OSS licenses from the license files of dependent libraries.

## :bookmark: Table of Contents

* [:runner: Usage](usage)
    * [`lioss`](usage/#lioss)
        * [Example](usage/#example)
    * [`mkliossdb`](usage/#mkliossdb)
* [:fork_and_knife: Install](install)
    * [:beer: Homebrew](install/#-homebrew)
    * [Go lang](install/#go-lang)
    * [:muscle: Build from source](install/#-build-from-source)
* [:package: LiossDB](liossdb)
* [:smile: About](about)
    * [:scroll: License](about/#-license)
    * [:man_office_worker: Developer :woman_office_worker:](about/#-developer-)
