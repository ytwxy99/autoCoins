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


    def getAll(self, tableName, coin):
        """get all records by specified database table

        :param tableName: the name of database table
        :param coin: the name of coin
        :return: all records by specified database table
        """
        c = self.conn.cursor()
        cursor = c.execute("SELECT * from %s where contract = '%s'" % (tableName, coin))

        return cursor