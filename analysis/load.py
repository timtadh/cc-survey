#!/usr/bin/env python

import os
import json


def load_answers(path):
    with open(path) as f:
        lines = [
            json.loads(line)
            for line in f
        ]
    d = dict()
    for line in lines:
        d[line['CloneID']] = line
    return d.values()

def prs(answers):
    return [
        line['ConditionalPr']
        for line in answers
    ]

def load_conditional_interval(base, clone_id):
    path = os.path.join(base, str(clone_id), 'conditional.interval.pr')
    with open(path) as f:
        return [
            float(line.strip())
            for line in f
        ]

def sel_prs(answers):
    return [
        line['SelectionPr']
        for line in answers
    ]

def edge_count(base, clone_id):
    path = os.path.join(base, str(clone_id), 'pattern.name')
    with open(path) as f:
        s = f.read().strip()
    size = s[:s.index('(')]
    E, V = size.split(':')
    return int(E)

def edge_counts(answers, base):
    return [
        edge_count(base, a['CloneExtID'])
        for a in answers
    ]

def duplicate(base, clone_id):
    path = os.path.join(base, str(clone_id), 'duplicates')
    with open(path) as f:
        s = f.read().strip()
    return int(s)

def duplicates(answers, base):
    return [
        duplicate(base, a['CloneExtID'])
        for a in answers
    ]

def prs_up(answers, base):
    return [
        load_conditional_interval(base, line['CloneExtID'])[1]
        for line in answers
    ]

def prs_down(answers, base):
    return [
        load_conditional_interval(base, line['CloneExtID'])[0]
        for line in answers
    ]

def question(answers, n, invert=False):
    if invert:
        return [
            abs(1 - line['Responses'][n]['Answer'])
            for line in answers
        ]
    else:
        return [
            line['Responses'][n]['Answer']
            for line in answers
        ]
