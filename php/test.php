<?php
include_once 'mtconfigserver.php';

$zk = new mtconfigserver();

$zk->init();
while(true){
    $value = $zk->get("inftest", "common", "test");
    
    //TODO your own business

    sleep(2);
}
?>
