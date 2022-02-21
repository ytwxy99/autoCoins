# -*- coding:utf8 -*-
import sys

import statsmodels.tsa.stattools as ts

from database import sqlite3
from utils import file
from utils import pandas


def initDatabase(dbPath):
    """init database object

    :param dbPath: the database path
    :return: the object of database
    """
    database = sqlite3.Database(sys.argv[1])
    database.initDb()
    return database


def pandasSeries(database, coins):
    """format pandas Series object

    :param database: the sqlite3 database object
    :param conis: all coins information
    """
    series = dict()
    for coin in coins:
        date = list()
        price = list()
        coinHistories = database.getAll("history_day", coin)
        for h in coinHistories:
            date.append(h[1])
            price.append(h[2])
        seriesObject = pandas.createPandasSeries(price, date)
        series[coin] = seriesObject

    return series


def getCointegration(coins, series):
    """get the cointegration of all coins"""
    storeCoints = dict()
    for coin in coins:
        if coin not in storeCoints:
            storeCoints[coin] = dict()
            storeCoints[coin]["coint"] = dict()
            storeCoints[coin]["coins"] = list()
        for c in coins:
            if c == coin:
                continue
            if c in storeCoints[coin]["coins"]:
                continue

            if len(series[coin]) == len(series[c]):
                cointRelation = coin + '-' + c
                coin_result = ts.coint(series[coin], series[c])
                storeCoints[coin]["coins"].append(c)
                storeCoints[coin]["coint"][cointRelation] = coin_result

    return storeCoints


def main():
    # ref: https://www.likecs.com/show-204274989.html
    try:
        cointResult = dict()
        database = initDatabase(sys.argv[1])
        coins = [coin.strip("\n") for coin in file.getFileContent(sys.argv[2])]
        series = pandasSeries(database, coins)
        storeCoints = getCointegration(coins, series)
        #NOTE(wangxiaoyu), cointegration judgment
        for coin in coins:
            for pair in storeCoints[coin]["coint"].keys():
                pValue = float(storeCoints[coin]["coint"][pair][1])
                if pValue <= 0.05 and pValue != 0.0:
                    cointResult[pair] = pValue

        coints = sorted(cointResult.items(), key=lambda x: x[1], reverse=False)
    except Exception as e:
        print(e)
    finally:
        print(coints)
        database.closeDB()


if __name__ == "__main__":
    main()