# -*- coding:utf8 -*-
import sqlite3


class Database(object):
    def __init__(self, dbPath):
        self.dbPath = dbPath
        self.conn = None

    def initDb(self):
        """get sqlite3 connection object

        :param dbPath: sqlite3 database file
        :return: database connection object
        """
        self.conn = sqlite3.connect(self.dbPath)

    def closeDB(self):
        """close sqlite3 connection

        :param dbConn:
        :return:
        """
        self.conn.close()

    def getHisotryDay(self, tableName, coin):
        """get all records by specified database table

        :param tableName: the name of database table
        :param coin: the name of coin
        :return: all records by specified database table
        """
        c = self.conn.cursor()
        cursor = c.execute("SELECT * from %s where contract = '%s'" % (tableName, coin))

        return cursor

    def getCointegration(self, tableName, pair):
        """get all records by specified database table

        :param tableName: the name of database table
        :param pair: the pair of cointegration
        :return: all records by specified database table
        """
        c = self.conn.cursor()
        cursor = c.execute("SELECT * from %s where pair = '%s'" % (tableName, pair))

        return cursor

    def insertCointegration(self, tableName, pair, pValue):
        """insert a cointegration data into Cointegration table

        :param tableName: the table name of database
        :param coint: the content of cointegration
        :return:
        """
        c = self.conn.cursor()
        c.execute("insert into %s (pair, pvalue) values ('%s', '%s')" % (tableName, pair, pValue))
        self.conn.commit()

    def updateCointegration(self, tableName, pair, pValue):
        """insert a cointegration data into Cointegration table

        :param tableName: the table name of database
        :param coint: the content of cointegration
        :return:
        """
        c = self.conn.cursor()
        c.execute("update %s set pvalue = '%s' where pair = '%s'" % (tableName, pValue, pair))
        self.conn.commit()