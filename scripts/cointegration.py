# -*- coding:utf8 -*-
import sys

from database import sqlite3
from utils import file


def initDatabase(dbPath):
    """init database object

    :param dbPath: the database path
    :return: the object of database
    """
    database = sqlite3.Database(sys.argv[1])
    database.initDb()
    return database


def main():
    try:
        database = initDatabase(sys.argv[1])
        coins = [coin.strip("\n") for coin in file.getFileContent(sys.argv[2])]
        for coin in coins:
            if coin == "ALPHR_USDT":
                coinHistories = database.getAll("history_day", coin)
                for h in coinHistories:
                    print(h)

    except Exception as e:
        print(e)
    finally:
        database.closeDB()


if __name__ == "__main__":
    main()