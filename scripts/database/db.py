# -*- coding:utf8 -*-
import sqlite3
import mysql.connector

class Database(object):
    def __init__(self, argv):
        self.conn = None
        self.dbType = argv[2]
        self.dbPath = argv[3]
        self.user = argv[4]
        self.password = argv[5]
        self.port = argv[6]
        self.host = argv[7]
        self.database = argv[8]

    def initDb(self):
        """get sqlite3 connection object

        :param dbPath: sqlite3 database file
        :return: database connection object
        """
        if self.dbType == "sqlite3":
            self.conn = sqlite3.connect(self.dbPath)
        elif self.dbType == "mysql":
            # we should set wait_time of mysql configure to 1h
            self.conn = mysql.connector.connect(
                host=self.host,
                user=self.user,
                password=self.password,
                database=self.database,
            )

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
        if self.dbType == "sqlite3":
            cursor = c.execute("SELECT * from %s where contract = '%s'" % (tableName, coin))
            return cursor
        elif self.dbType == "mysql":
            c.execute("SELECT * from %s where contract = '%s'" % (tableName, coin))
            return c.fetchall()

    def getCointegration(self, tableName, pair):
        """get all records by specified database table

        :param tableName: the name of database table
        :param pair: the pair of cointegration
        :return: all records by specified database table
        """
        c = self.conn.cursor()
        if self.dbType == "sqlite3":
            cursor = c.execute("SELECT * from %s where pair = '%s'" % (tableName, pair))
            return cursor
        elif self.dbType == "mysql":
            c.execute("SELECT * from %s where pair = '%s'" % (tableName, pair))
            return c

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