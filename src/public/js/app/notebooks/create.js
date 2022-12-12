
$(document).on('show.bs.modal','#notebook-modal', function (e) {
  const target = $(e.relatedTarget);
  $('#notebook-menu-id').val(target.data('menu'));
  $('#notebook-menu').text(target.closest('.card').find('.card-header button').text());
});

$(document).on("click", "#notebook-modal .btn-primary", function() {
  $('#notebook-modal').modal('hide');

  const params = {menu: $('#notebook-menu-id').val()};
  const option = {headers: {Authorization: `Bearer ${idToken}`}};
  axios
    .post(window.apiBaseURL+'/api/notebooks', params, option)
    .then(_ => location.href = '/notebooks')
    .catch(error => console.log(error));
});
