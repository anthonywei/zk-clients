<?php
/**
 * mtconfigserver for php 
 *
 * 2014-06-30.
 */
/**
 * the url of the config server
 */
defined("MTCONFIG_SERVER") or define("MTCONFIG_SERVER", "http://xxx.xxx.com/"); //管理中心默认的url
defined("MTCONFIG_PATH") or define("MTCONFIG_PATH", "/config");     //配置默认的根路径

class mtconfigserver
{
    private $config_server = MTCONFIG_SERVER;
	private $zookeeper;
    private $cache_config;      //config array
    private $zk_list;

	public function __construct($config_server_ = MTCONFIG_SERVER) {
        $this->config_server = $config_server_;

        $this->init();
	}

    public function init() {
        $curl = curl_init();
        curl_setopt ($curl, CURLOPT_URL, $this->config_server . "/api/zkserverlist");
        curl_setopt ($curl, CURLOPT_HEADER, 0);
        curl_setopt ($curl, CURLOPT_RETURNTRANSFER, 1);
        curl_setopt ($curl, CURLOPT_TIMEOUT,5);

        $get_content = curl_exec($curl);
        curl_close ($curl);

        $this->zk_list = $get_content;

        $this->zookeeper = new Zookeeper($this->zk_list);
    }

    public function get($domain, $node, $key) {
        $path = MTCONFIG_PATH . "/". $domain . "/" . $node;


        //get from cache
        if(isset($this->cache_config[$path][$key]))
            return $this->cache_config[$path][$key];

        //get from zk
        $value = $this->_get($path);

        return $this->cache_config[$path][$key];
    }

	private function _get($path) {
		if (!$this->zookeeper->exists($path)) {
			return null;
		}

        /* delete old configs*/
        if(isset($this->cache_config[$path]))
            unset($this->cache_config[$path]);

        /* try to get new configs*/
		$content = $this->zookeeper->get($path, array($this, 'onCallBack'));
        
        /* try to parse return contents*/
        $content = trim($content);
        $content_array = explode("\n\n", $content);
        for($i = 0; $i < count($content_array); $i ++) {
            $item_array = explode("=", $content_array[$i]);
            if(count($item_array) == 2) {
                $key = trim($item_array[0]);
                $value = trim($item_array[1]);

                /* set it to cache */
                $this->cache_config[$path][$key] = $value;
            }
        }

        return $content;
	}


    public function onCallBack($event_type, $stat, $path) {
        /* get the config again */
        $this->_get($path); 
    }
}



