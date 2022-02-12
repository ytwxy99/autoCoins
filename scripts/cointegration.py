# -*- coding:utf8 -*-
import sys

from database import sqlite3
from utils import file


def main():
    conn = sqlite3.getDb(sys.argv[1])
    coins = [coin.strip("\n") for coin in file.getFileContent(sys.argv[2])]
    for coin in coins:
        coinHistories = sqlite3.getAll("history_day", coin, conn)
        print(coinHistories)

    conn.close()


if __name__ == "__main__":
    main()