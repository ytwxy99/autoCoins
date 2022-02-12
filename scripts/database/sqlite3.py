# -*- coding:utf8 -*-
import sqlite3

def getDb(dbPath):
    """get sqlite3 connection object

    :param dbPath: sqlite3 database file
    :return: database connection object
    """
    return sqlite3.connect(dbPath)


def closeDB(dbConn):
    """close sqlite3 connection

    :param dbConn:
    :return:
    """
    dbConn.close()


def getAll(tableName, coin,  dbConn):
    """get all records by specified database table

    :param tableName: the name of database table
    :param coin: the name of coin
    :param dbConn: sqlite3 connection object
    :return: all records by specified database table
    """
    c = dbConn.cursor()

    cursor = c.execute("SELECT * from %s where contract = '%s'" % (tableName, coin))
    for row in cursor:
        print(row[0], row[1], row[2])

    return cursor