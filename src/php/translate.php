<?php
function openai_client($url, $token, $sys_msg, $desc) {

}
function get_commit_by_open_ai($desc) {

}
if ($_SERVER["REQUEST_METHOD"] === 'POST') {
    //收集 form 的数据
    $description = $_POST["description"];

    // 你的 OpenAI secret key
    $secret_key = 'sk-rLEbROmxDX3MfRq2155e4138BbAb48699e8c6203D492FaE1';

    // get open ai token by env
    $secret_key = getenv('OPENAI_SECRET_KEY');

    // 向 OpenAI API 发起请求
    $host = 'https://one-api.skybyte.me';
    $url = $host.'/v1/chat/completions';
    $system_msg = "根据我的描述，为我生成简短的英文 git commit \
                    消息，格式为 “<type>(<scope>): <subject>”，\
                    其中type用于说明 commit 的类别，只允许使用\
                    feat、fix、docs、style、refactor、test、chore，\
                    scope用于说明 commit 影响的范围，\
                    subject是 commit 目的的简短描述，不超过50个字符。\
                    如果输入有误或空白，请输出 Error.";
    $desc = $description;
    $data = json_encode([
        "model" => 'gpt-3.5-turbo',
        "messages" => [
            [
                "role" => "system",
                "content" => $system_msg
            ],
            [
                "role" => "user",
                "content" => $desc,
            ]
        ],
        "temperature" => 1.0,
        "max_tokens" => 256,
        "top_p" => 1.0,
        "frequency_penalty" => 0.0,
        "presence_penalty" => 0.0
    ]);

    $options = array(
        'http' => array(
            'header'  => "Content-type: application/json\r\nAuthorization: Bearer ".$secret_key,
            'method'  => 'POST',
            'content' => $data
        )
    );

    $context  = stream_context_create($options);

    // 发送请求
    $result = file_get_contents($url, false, $context);

    // 检查是否有错误
    if($result === FALSE){
        die('Error occurred!');
    }

    // 解析响应
    $responseData = json_decode($result, TRUE);

    // 输出响应
    echo $responseData["choices"][0]["message"]["content"];
}

