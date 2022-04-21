# -*- coding:utf8 -*-
import sys

import statsmodels.tsa.stattools as ts

from database import db
from utils import file
from utils import pandas

COINTEGRATION_DB = "cointegration"

def initDatabase(argv):
    """init database object

    :param dbPath: argv
    :return: the object of database
    """
    database = db.Database(sys.argv)
    database.initDb()
    return database


def pandasSeries(database, coins):
    """format pandas Series object

    :param database: the database object
    :param conis: all coins information
    """
    series = dict()
    for coin in coins:
        date = list()
        price = list()
        coinHistories = database.getHisotryDay("history_day", coin)
        for h in coinHistories:
            date.append(h[1])
            price.append(h[2])
        seriesObject = pandas.createPandasSeries(price, date)
        series[coin] = seriesObject

    return series


def getCointegration(coins, series):
    """get the cointegration of all coins

    :param coins: all coins information
    :param series: coins series
    :return: the cointegration of coins
    """
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
                #todo(wangxiaoyu), find out why an error is reported
                # when the length of series[coins] and series[c] is 1
                if len(series[coin]) < 900:
                    continue

                cointRelation = coin + '-' + c
                coin_result = ts.coint(series[coin], series[c])
                storeCoints[coin]["coins"].append(c)
                storeCoints[coin]["coint"][cointRelation] = coin_result

    return storeCoints


def main():
    # ref: https://www.likecs.com/show-204274989.html
    try:
        cointResult = dict()
        database = initDatabase(sys.argv)

        coins = [coin.strip("\n") for coin in file.getFileContent(sys.argv[1])]
        series = pandasSeries(database, coins)
        storeCoints = getCointegration(coins, series)
        #NOTE(wangxiaoyu), cointegration judgment
        for coin in coins:
            for pair in storeCoints[coin]["coint"].keys():
                pValue = float(storeCoints[coin]["coint"][pair][1])
                if pValue <= 0.05 and pValue != 0.0:
                    cointResult[pair] = pValue

        coints = sorted(cointResult.items(), key=lambda x: x[1], reverse=False)

        database = initDatabase(sys.argv)
        for coint in coints:
            cointP = database.getCointegration(COINTEGRATION_DB, coint[0]).fetchall()
            if len(cointP) > 0:
                database.updateCointegration(COINTEGRATION_DB, coint[0], coint[1])
            else:
                database.insertCointegration(COINTEGRATION_DB, coint[0], coint[1])

    except Exception as e:
        print(e)
    finally:
        database.closeDB()


if __name__ == "__main__":
    main()