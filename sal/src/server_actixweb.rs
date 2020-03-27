use crate::types;
use actix_web;

fn error_json<T: std::string::ToString>(e: T) -> String {
    format!("{{\"error\":\"{}\"}}", e.to_string())
}

pub mod backend {
    use super::*;

    pub fn routes<T>(cfg: &mut actix_web::web::ServiceConfig)
    where
        T: types::Backend + std::clone::Clone + 'static,
    {
        cfg.service(actix_web::web::resource("/oto/Backend.Run").to(run::<T>))
            .service(actix_web::web::resource("/oto/Backend.Test").to(test::<T>))
            .service(actix_web::web::resource("/oto/Backend.Visualize").to(visualize::<T>));
    }

    // address: e.g. "127.0.0.1:8080"
    pub async fn main<T>(svc: T, address: &str) -> std::io::Result<()>
    where
        T: types::Backend + std::clone::Clone + 'static,
    {
        actix_web::HttpServer::new(move || {
            actix_web::App::new()
                .data(svc.clone())
                .configure(routes::<T>)
        })
        .bind(address)?
        .run()
        .await
    }

    pub async fn run<'a, T>(
        svc: actix_web::web::Data<T>,
        req: actix_web::web::Json<types::RunRequest>,
    ) -> impl actix_web::Responder
    where
        T: types::Backend + 'a,
    {
        let (status, body) = match svc.run(req.into_inner()).await {
            Ok(res) => match serde_json::to_string(&res) {
                Ok(body) => (actix_web::http::StatusCode::OK, body),
                Err(e) => (
                    actix_web::http::StatusCode::INTERNAL_SERVER_ERROR,
                    error_json(format!("error serializing response: {:?}", e)),
                ),
            },
            Err(e) => (
                actix_web::http::StatusCode::INTERNAL_SERVER_ERROR,
                error_json(e),
            ),
        };
        actix_web::web::HttpResponse::build(status)
            .content_type("application/json")
            .body(body)
    }

    pub async fn test<'a, T>(
        svc: actix_web::web::Data<T>,
        req: actix_web::web::Json<types::TestRequest>,
    ) -> impl actix_web::Responder
    where
        T: types::Backend + 'a,
    {
        let (status, body) = match svc.test(req.into_inner()).await {
            Ok(res) => match serde_json::to_string(&res) {
                Ok(body) => (actix_web::http::StatusCode::OK, body),
                Err(e) => (
                    actix_web::http::StatusCode::INTERNAL_SERVER_ERROR,
                    error_json(format!("error serializing response: {:?}", e)),
                ),
            },
            Err(e) => (
                actix_web::http::StatusCode::INTERNAL_SERVER_ERROR,
                error_json(e),
            ),
        };
        actix_web::web::HttpResponse::build(status)
            .content_type("application/json")
            .body(body)
    }

    pub async fn visualize<'a, T>(
        svc: actix_web::web::Data<T>,
        req: actix_web::web::Json<types::VisualizeRequest>,
    ) -> impl actix_web::Responder
    where
        T: types::Backend + 'a,
    {
        let (status, body) = match svc.visualize(req.into_inner()).await {
            Ok(res) => match serde_json::to_string(&res) {
                Ok(body) => (actix_web::http::StatusCode::OK, body),
                Err(e) => (
                    actix_web::http::StatusCode::INTERNAL_SERVER_ERROR,
                    error_json(format!("error serializing response: {:?}", e)),
                ),
            },
            Err(e) => (
                actix_web::http::StatusCode::INTERNAL_SERVER_ERROR,
                error_json(e),
            ),
        };
        actix_web::web::HttpResponse::build(status)
            .content_type("application/json")
            .body(body)
    }
}

#[cfg(test)]
mod tests {

    mod backend {
        use super::super::*;
        use actix_web::test;

        async fn index() -> impl actix_web::Responder {
            String::from("Hello, world!")
        }

        #[actix_rt::test]
        async fn run_ok() {
            let endpoint = "/oto/Backend.Run";
            let body = types::RunRequest::new();
            let req = test::TestRequest::post()
                .uri(endpoint)
                .set_json(&body)
                .to_request();
            let svc = types::mock::MockBackend::new();
            let mut app = test::init_service(
                actix_web::App::new()
                    .data(svc.clone())
                    .configure(super::super::backend::routes::<types::mock::MockBackend>),
            )
            .await;
            let resp = test::call_service(&mut app, req).await;
            assert_eq!(resp.status(), actix_web::http::StatusCode::OK);
            let result = test::read_body(resp).await;
            let obj = types::RunResponse::new();
            let expected = serde_json::to_string(&obj).unwrap();
            assert_eq!(result, expected);
        }

        #[actix_rt::test]
        async fn run_error() {
            let endpoint = "/oto/Backend.Run";
            let body = types::RunRequest::new();
            let req = test::TestRequest::post()
                .uri(endpoint)
                .set_json(&body)
                .to_request();
            let svc = types::mock::MockBackend::error("Hello from Backend.Run!");
            let mut app = test::init_service(
                actix_web::App::new()
                    .data(svc.clone())
                    .configure(super::super::backend::routes::<types::mock::MockBackend>),
            )
            .await;
            let resp = test::call_service(&mut app, req).await;
            assert_eq!(
                resp.status(),
                actix_web::http::StatusCode::INTERNAL_SERVER_ERROR
            );
            let result = test::read_body(resp).await;
            let expected = error_json("Hello from Backend.Run!");
            assert_eq!(result, expected);
        }

        #[actix_rt::test]
        async fn test_ok() {
            let endpoint = "/oto/Backend.Test";
            let body = types::TestRequest::new();
            let req = test::TestRequest::post()
                .uri(endpoint)
                .set_json(&body)
                .to_request();
            let svc = types::mock::MockBackend::new();
            let mut app = test::init_service(
                actix_web::App::new()
                    .data(svc.clone())
                    .configure(super::super::backend::routes::<types::mock::MockBackend>),
            )
            .await;
            let resp = test::call_service(&mut app, req).await;
            assert_eq!(resp.status(), actix_web::http::StatusCode::OK);
            let result = test::read_body(resp).await;
            let obj = types::TestResponse::new();
            let expected = serde_json::to_string(&obj).unwrap();
            assert_eq!(result, expected);
        }

        #[actix_rt::test]
        async fn test_error() {
            let endpoint = "/oto/Backend.Test";
            let body = types::TestRequest::new();
            let req = test::TestRequest::post()
                .uri(endpoint)
                .set_json(&body)
                .to_request();
            let svc = types::mock::MockBackend::error("Hello from Backend.Test!");
            let mut app = test::init_service(
                actix_web::App::new()
                    .data(svc.clone())
                    .configure(super::super::backend::routes::<types::mock::MockBackend>),
            )
            .await;
            let resp = test::call_service(&mut app, req).await;
            assert_eq!(
                resp.status(),
                actix_web::http::StatusCode::INTERNAL_SERVER_ERROR
            );
            let result = test::read_body(resp).await;
            let expected = error_json("Hello from Backend.Test!");
            assert_eq!(result, expected);
        }

        #[actix_rt::test]
        async fn visualize_ok() {
            let endpoint = "/oto/Backend.Visualize";
            let body = types::VisualizeRequest::new();
            let req = test::TestRequest::post()
                .uri(endpoint)
                .set_json(&body)
                .to_request();
            let svc = types::mock::MockBackend::new();
            let mut app = test::init_service(
                actix_web::App::new()
                    .data(svc.clone())
                    .configure(super::super::backend::routes::<types::mock::MockBackend>),
            )
            .await;
            let resp = test::call_service(&mut app, req).await;
            assert_eq!(resp.status(), actix_web::http::StatusCode::OK);
            let result = test::read_body(resp).await;
            let obj = types::VisualizeResponse::new();
            let expected = serde_json::to_string(&obj).unwrap();
            assert_eq!(result, expected);
        }

        #[actix_rt::test]
        async fn visualize_error() {
            let endpoint = "/oto/Backend.Visualize";
            let body = types::VisualizeRequest::new();
            let req = test::TestRequest::post()
                .uri(endpoint)
                .set_json(&body)
                .to_request();
            let svc = types::mock::MockBackend::error("Hello from Backend.Visualize!");
            let mut app = test::init_service(
                actix_web::App::new()
                    .data(svc.clone())
                    .configure(super::super::backend::routes::<types::mock::MockBackend>),
            )
            .await;
            let resp = test::call_service(&mut app, req).await;
            assert_eq!(
                resp.status(),
                actix_web::http::StatusCode::INTERNAL_SERVER_ERROR
            );
            let result = test::read_body(resp).await;
            let expected = error_json("Hello from Backend.Visualize!");
            assert_eq!(result, expected);
        }
    }
}
