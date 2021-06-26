<?php
//Read public key and generate unique id
$pubKey = file_get_contents("key.pub");
$id = uniqid();

//Write into base
$current = file_get_contents("base.key");
$current .= $id . ":\n";
file_put_contents("base.key", $current);

//Send data in json format
$data = array('id' => $id, 'p licKey' => $pubKey);
echo json_encode($data);
?>
