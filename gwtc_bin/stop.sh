#!/bin/bash
ps -ef | grep gwtc | awk '{ print $2 }' | sudo xargs kill -9
