# coding: utf-8
# auth weishouyang@meituan.com 2014-06-30
# mtconfigserver

import sys
import zookeeper, time, threading, urllib

class mtconfigserver():
    def __init__(self, config_server_="http://xxx.xxx.com/"):
        self.config_server = config_server_
        self.cache_config = []  
        self.zk_path = "/config"
        self.zklist = self.getZkServer()
        zookeeper.set_debug_level(zookeeper.LOG_LEVEL_ERROR)
        self.zookeeper = zookeeper.init(self.zklist)

    def getZkServer(self):
        url = urllib.urlopen(self.config_server + "/api/zkserverlist")
        data = url.read()
        return data
    
    def get(self, domain, node, key):
        path = self.zk_path + "/" + domain + "/" + node
        value = self._find_in_cache(path, key)
        if value != "":
            return value

        self._get_from_zk(path)
        return self._find_in_cache(path, key)


    def _find_in_cache(self, path, key):
        for cs in self.cache_config:
            if cs.path==path:
                for ci in cs.items:
                    if ci.key==key:
                        return ci.value
        return ""

    def _delete_from_cache(self, path):
        i = 0
        while i<len(self.cache_config):
            if self.cache_config[i].path == path:
                del self.cache_config[i]
                break
            i=i+1 

    def _get_from_zk(self, path):
        self._delete_from_cache(path)
        content = zookeeper.get(self.zookeeper, path, self._onCallBack)
        self._parse(path, content[0])

    def _parse(self, path, config_data):
        strip_data = config_data.strip()
        data_array = strip_data.split('\n\n')
        cs = config_set(path)
        i = 0
        while i<len(data_array):
            item_array = data_array[i].split('=')
            if len(item_array)==2:
                key = item_array[0].strip()
                value = item_array[1].strip()
                ci = config_item(key, value)
                cs.items.append(ci)
            i=i+1

        self.cache_config.append(cs)


    def _onCallBack(self, handler, event_type, stat, path):
        self._get_from_zk(path)

class config_set():
    def __init__(self, path):
        self.path = path
        self.items = []

class config_item():
    def __init__(self, key, value):
        self.key = key
        self.value = value



