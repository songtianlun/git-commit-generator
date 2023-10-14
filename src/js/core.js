
$(document).ready(function(){
    $('#myForm').on('submit', function(e) {
        e.preventDefault();
        $.ajax({
            url: 'src/php/translate.php', // 替换为实际的 PHP API URL
            method: 'post',
            data: $('#myForm').serialize(),
            success: function(response) {
                $('#result').html(response); // 将结果填入 #result 元素
            }
        });
    });

    var input = document.getElementById("result");
    var button = document.getElementById("copyButton");

    button.addEventListener("click", function() {
        input.select();
        document.execCommand("copy");
    });
});