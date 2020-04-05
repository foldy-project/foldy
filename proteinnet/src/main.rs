use actix_web;


pub fn routes(cfg: &mut actix_web::web::ServiceConfig)  {
    //cfg.service(actix_web::web::resource("/oto/Driver.Associate").to(associate::<T>));
}

#[actix_rt::main]
async fn main() -> std::io::Result<()> {
    let address = "0.0.0.0:8080";
    actix_web::HttpServer::new(move || {
        actix_web::App::new()
            .wrap(actix_web::middleware::Logger::default())
            .configure(routes)
    })
        .bind(address)?
        .run()
        .await
}

