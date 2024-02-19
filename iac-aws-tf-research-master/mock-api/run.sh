#!/bin/bash
pipenv install
pipenv run moto_server -p 3000
