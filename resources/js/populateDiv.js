
//every time page reloads(or first time on home page), show all results
$( document ).ready(function() {
    updateResults('all');
});
$( '#search_button' ).click(function(){
    updateResults('all');
});
//These functions are when page is not reloading, static page
function showFilter(data){         
    var fields= JSON.parse(data);      
    var box="";
    $.each(fields, function(field , obj) {
        box += '<div class="row">';
        box += '<div class="filter-first">';
        box += '<a href="#add" class="filter-by-type">'+ field + '</a></div>';
        $.each(obj, function(subf , count) {
        //if count ==0 then dont display checkbox
            if(count >0){
                box +='<div class="checkbox ">'; 
                box +='<label id="subfields">';
                box +='<input type="checkbox" name='+ field +' value='+ subf +'>'+ subf + "("+count+")";
                box +='</label></div>';
            }
        });
        box += '</div>';
        });
        document.getElementById("filterbox").innerHTML= box; 
}      
//These functions are when page is not reloading, static page
function showResults(data){
    var animals= JSON.parse(data);
    var box="";
    $.each(animals, function(k , v) {
        var href= "../animallistings/view/" + v[0] +"/"+v[1];
        box += "<div id='arect' class='col-md-4 portfolio-item'>";
        box += '<a class="thumbnail" href='+ href +'>';
        box += '<img class="img-responsive" width="260" height="180" src= "http://sfsuswe.com/~s15g12/Framework/'+ v["image"]["image_path"]+'" alt="No Image">';
        box += '<span class="animal">';
        box += '<div class="caption">';
        box += v[1] + '</div></span></a></div>';
    });
    document.getElementById("box").innerHTML= box; 
}
//Store all checked options in opts array
function getFilterOptions(){
    var opts = {};
    $checkboxes.each(function(){
        if(this.checked){
            alert(this.name);
            opts[this.name]= this.value;    
        }
    });
    opts = JSON.stringify(opts);
    return opts;
}
function updateResults(keyword,opts){
    //post to index, get results and showResults
    $.ajax({
        type: "POST",
        url: "searching",
        data: {search: keyword, filterOpts: opts},
        dataType : 'html',
        success: function(data) {        
            var matches = data.substr(0, data.indexOf('{"Animals":')); 
            // var fil = data.substr(data.indexOf('{"Animals":'), data.indexOf("3}}")); 
            var filter_data = data.match('{"Animals":(.*)}}');
            //insert the formatted div in the show results div
            showResults(matches);
            //update Filter labels for the next query
            showFilter(filter_data[0]);
        },
        //error: function(xhr, textStatus, errorThrown) {
            //alert("Fail");
       // }
    });
}

