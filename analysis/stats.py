#!/usr/bin/env python

import math

import numpy as np
from scipy import stats


class SampleProbabilities(object):

    def __init__(self, prs, size, pr_subpop):
        self.n = size
        self.pr_s = np.mean(pr_subpop)
        self._prs = prs
        self._pis = None
        self._jpis = None

    @property
    def prs(self):
        return self._prs

    @property
    def pis(self):
        if self._pis is None:
            self._pis = self.c_pis(self.prs, self.n)
        return self._pis

    @property
    def jpis(self):
        if self._jpis is None:
            self._jpis = self.c_jpis(self.prs, self.pis, self.n)
        return self._jpis

    def pi(self, i, prs, n):
        return 1.0 - ((1.0 - prs[i])**n)

    def jpi(self, i, j, prs, pis, n):
        return pis[i] + pis[j] - (1.0 - (1.0 - prs[i] - prs[j])**n)

    def c_pis(self, prs, n):
        return [
            self.pi(i, prs, n)
            for i in xrange(len(prs))
        ]

    def c_jpis(self, prs, pis, n):
        return [
            [ self.jpi(i, j, prs, pis, n) for j in xrange(len(prs)) ]
            for i in xrange(len(prs))
        ]

class Estimators(object):

    def __init__(self, ys, sample_probabilities, confidence=.95):
        self.ys = ys
        self.p = sample_probabilities
        self.a = stats.t.ppf(confidence, self.p.n)
        self._tau_hat = None
        self._var_tau_hat = None
        self._n_hat = None
        self._var_n_hat = None
        self._mu_hat = None
        self._var_mu_hat = None


    @property
    def tau_hat(self):
        if self._tau_hat is None:
            self._tau_hat = self.c_tau_hat(self.ys, self.p.pis)
        return self._tau_hat

    @property
    def var_tau_hat(self):
        if self._var_tau_hat is None:
            self._var_tau_hat = self.c_var_tau_hat(
                    self.ys, self.p.pis, self.p.jpis)
        return self._var_tau_hat

    @property
    def std_tau_hat(self):
        return math.sqrt(self.var_tau_hat)

    @property
    def interval_tau_hat(self):
        i = self.a * self.std_tau_hat
        return (self.tau_hat + i, self.tau_hat - i)

    @property
    def n_hat(self):
        if self._n_hat is None:
            self._n_hat = self.c_n_hat(self.p.pis)
        return self._n_hat

    @property
    def var_n_hat(self):
        if self._var_n_hat is None:
            self._var_n_hat = self.c_var_n_hat(self.p.pis, self.p.jpis)
        return self._var_n_hat

    @property
    def std_n_hat(self):
        return math.sqrt(self.var_n_hat)

    @property
    def interval_n_hat(self):
        i = self.a * self.std_n_hat
        return (self.n_hat + i, self.n_hat - i)

    @property
    def mu_hat(self):
        if self._mu_hat is None:
            self._mu_hat = self.c_mu_hat(
                self.n_hat, self.ys, self.p.pis)
        return self._mu_hat

    @property
    def var_mu_hat(self):
        if self._var_mu_hat is None:
            self._var_mu_hat = self.c_var_mu_hat(
                self.n_hat, self.mu_hat,
                self.ys, self.p.pis, self.p.jpis)
            ## self._var_mu_hat = self.c_var_p_hat(
            ##      self.p.n, self.n_hat, self.mu_hat)
        return self._var_mu_hat

    @property
    def std_mu_hat(self):
        return math.sqrt(self.var_mu_hat)

    @property
    def interval_mu_hat(self):
        i = self.a * self.var_mu_hat
        return (self.mu_hat + i, self.mu_hat - i)

class HT_Estimators(Estimators):

    def c_tau_hat(self, ys, pis):
        return sum(ys[i]/pis[i] for i in xrange(len(pis)))

    def c_tau_hat(self, ys, pis):
        return sum(ys[i]/pis[i] for i in xrange(len(pis)))

    def c_var_tau_hat(self, ys, pis, jpis):
        a = sum(
            (1.0/(pis[i]**2) - (1.0/pis[i]))
            for i in xrange(len(pis)))
        b = sum(
            ys[i]*ys[j]*((1.0/(pis[i]*pis[j])) - (1.0/jpis[i][j]))
            for j in xrange(len(pis))
            for i in xrange(len(pis))
            if i != j
        )
        return a + 2.0*b

    def c_n_hat(self, pis):
        return self.c_tau_hat([1 for _ in xrange(len(pis))], pis)

    def c_var_n_hat(self, pis, jpis):
        return self.c_var_tau_hat([1 for _ in xrange(len(pis))], pis, jpis)

    def c_mu_hat(self, n_hat, ys, pis):
        N = self.p.n/self.p.pr_s
        return (1.0/(n_hat))*self.c_tau_hat(ys, pis)

    def c_var_mu_hat(self, n_hat, mu, ys, pis, jpis):
        a = sum(
                (1.0 - pis[i])/(pis[i]**2)
                for i in xrange(len(pis)))
        b = sum(
            (((jpis[i][j] - pis[i]*pis[j])/(pis[i]*pis[j])) * (((ys[i] - mu)*(ys[j] - mu))/(jpis[i][j])))
            for j in xrange(len(pis))
            for i in xrange(len(pis))
            if i != j
        )
        scale = (1.0/(n_hat**2))
        print 'scale', scale
        print 'a', a
        print 'b', b
        return scale*(a + b)

    ## This is the estimate of the variance of the estimate of the proportion
    ## as given by thompson on page 58. it seems to be work better?
    def c_var_p_hat(self, n, n_hat, p_hat):
        return ((n_hat - n)/(n_hat))*((p_hat*(1 - p_hat))/(n - 1))

class P_Estimators(HT_Estimators):
    ''' proportion estimators uses HT for tau'''

    def c_mu_hat(self, n_hat, ys, pis):
        N = self.p.n/self.p.pr_s
        print N
        print self.p.n
        print sum(ys[i]*self.p.prs[i] for i in xrange(len(ys)))
        a = self.p.n*sum(ys[i]*self.p.prs[i] for i in xrange(len(ys)))
        return a

    def c_var_mu_hat(self, n_hat, mu, ys, pis, jpis):
        return mu - (mu**2)
