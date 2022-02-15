# -*- coding:utf8 -*-
import pandas as pd


def createPandasSeries(data, index=None):
    """create the struct of pandas series object

    :param data: the data of pandas series
    :param index: the index of pandas index
    :return: pandas series object
    """
    if index is None:
        return pd.Series(data)
    else:
        return pd.Series(data, index)