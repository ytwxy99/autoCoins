# -*- coding:utf8 -*-

def getFileContent(path):
    """get specified file content

    :param path: file path
    :return: the content of the file
    """
    return open(path, "r").readlines()