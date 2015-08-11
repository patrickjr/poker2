
function toggleThis(divID){
        $(divID).toggle();
}
function toggleThisSlow(divID){
        $(divID).animate({width:"toggle", height: "toggle"},325);
}
/*function toggleThree(divClassLink, divClassContent){
    var $this = $(this);
    var target = $this.data(divClassContent);
    $(divClassLink).not($this).each(function(){
        var $other = $(this);
        var otherTarget = $other.data(divClassContent);
       $(otherTarget).hide();        
    });
    $(target).animate({width:"toggle", height: "toggle"},200);

}


function toggleAdmin(divClass){
    var $this   = $(this);
    var target = $this.data('content');
    $('div'+divClass).not($this).each(function(){
        var $other = $(this);
        var otherTarget = $other.data('content');
        $(otherTarget).hide();        
    });
    $(target).toggle();
}*/
