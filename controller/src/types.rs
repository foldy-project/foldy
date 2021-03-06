// Code generated by oto; DO NOT EDIT.
use async_trait::async_trait;
use serde::{Deserialize, Serialize};
use serde_json::{Map, Value};

#[async_trait]
pub trait Backend: std::marker::Send + std::marker::Sync {
    async fn run(&self, req: RunRequest) -> Result<RunResponse, String>;
    async fn test(&self, req: TestRequest) -> Result<TestResponse, String>;
    async fn visualize(&self, req: VisualizeRequest) -> Result<VisualizeResponse, String>;
}

#[derive(Serialize, Deserialize)]
pub struct Atom {
    pub id: String,
    pub element: String,
    pub x: f32,
    pub z: f32,
    pub y: f32,
}

impl Atom {
    pub fn new() -> Self {
        Atom {
            id: String::new(),
            element: String::new(),
            x: 0.0,
            z: 0.0,
            y: 0.0,
        }
    }

    pub fn make(id: String, element: String, x: f32, z: f32, y: f32) -> Self {
        Atom {
            id,

            element,

            x,

            z,

            y,
        }
    }
}

impl Default for Atom {
    fn default() -> Atom {
        Atom::new()
    }
}

#[derive(Serialize, Deserialize)]
pub struct RunRequest {
    pub id: String,
    pub backend: String,
    pub input: String,
    pub config: std::collections::HashMap<String, Value>,
    pub foo: std::collections::HashMap<i64, String>,
}

impl RunRequest {
    pub fn new() -> Self {
        RunRequest {
            id: String::new(),
            backend: String::new(),
            input: String::new(),
            config: std::collections::HashMap::<String, Value>::new(),
            foo: std::collections::HashMap::<i64, String>::new(),
        }
    }

    pub fn make(
        id: String,
        backend: String,
        input: String,
        config: std::collections::HashMap<String, Value>,
        foo: std::collections::HashMap<i64, String>,
    ) -> Self {
        RunRequest {
            id,

            backend,

            input,

            config,

            foo,
        }
    }
}

impl Default for RunRequest {
    fn default() -> RunRequest {
        RunRequest::new()
    }
}

#[derive(Serialize, Deserialize)]
pub struct RunResponse {
    error: Option<String>,
}

impl RunResponse {
    pub fn new() -> Self {
        RunResponse { error: None }
    }

    pub fn make() -> Self {
        RunResponse { error: None }
    }

    pub fn error(msg: String) -> Self {
        RunResponse { error: Some(msg) }
    }

    pub(crate) fn is_error(&self) -> bool {
        self.error != None
    }

    pub(crate) fn get_error(&self) -> Option<&String> {
        self.error.as_ref()
    }

    pub(crate) fn take_error(&mut self) -> Option<String> {
        self.error.take()
    }
}

impl Default for RunResponse {
    fn default() -> RunResponse {
        RunResponse::new()
    }
}

#[derive(Serialize, Deserialize)]
pub struct TestRequest {}

impl TestRequest {
    pub fn new() -> Self {
        TestRequest {}
    }

    pub fn make() -> Self {
        TestRequest {}
    }
}

impl Default for TestRequest {
    fn default() -> TestRequest {
        TestRequest::new()
    }
}

#[derive(Serialize, Deserialize)]
pub struct TestResponse {
    error: Option<String>,
}

impl TestResponse {
    pub fn new() -> Self {
        TestResponse { error: None }
    }

    pub fn make() -> Self {
        TestResponse { error: None }
    }

    pub fn error(msg: String) -> Self {
        TestResponse { error: Some(msg) }
    }

    pub(crate) fn is_error(&self) -> bool {
        self.error != None
    }

    pub(crate) fn get_error(&self) -> Option<&String> {
        self.error.as_ref()
    }

    pub(crate) fn take_error(&mut self) -> Option<String> {
        self.error.take()
    }
}

impl Default for TestResponse {
    fn default() -> TestResponse {
        TestResponse::new()
    }
}

#[derive(Serialize, Deserialize)]
pub struct VisualizeRequest {
    pub f_p_s: i64,
    pub num_frames: i64,
}

impl VisualizeRequest {
    pub fn new() -> Self {
        VisualizeRequest {
            f_p_s: 0,
            num_frames: 0,
        }
    }

    pub fn make(f_p_s: i64, num_frames: i64) -> Self {
        VisualizeRequest { f_p_s, num_frames }
    }
}

impl Default for VisualizeRequest {
    fn default() -> VisualizeRequest {
        VisualizeRequest::new()
    }
}

#[derive(Serialize, Deserialize)]
pub struct VisualizeResponse {
    pub resource: String,
    error: Option<String>,
}

impl VisualizeResponse {
    pub fn new() -> Self {
        VisualizeResponse {
            resource: String::new(),
            error: None,
        }
    }

    pub fn make(resource: String) -> Self {
        VisualizeResponse {
            resource,

            error: None,
        }
    }

    pub fn error(msg: String) -> Self {
        VisualizeResponse {
            error: Some(msg),

            resource: Default::default(),
        }
    }

    pub(crate) fn is_error(&self) -> bool {
        self.error != None
    }

    pub(crate) fn get_error(&self) -> Option<&String> {
        self.error.as_ref()
    }

    pub(crate) fn take_error(&mut self) -> Option<String> {
        self.error.take()
    }
}

impl Default for VisualizeResponse {
    fn default() -> VisualizeResponse {
        VisualizeResponse::new()
    }
}

#[derive(Serialize, Deserialize)]
pub struct Residue {
    pub atoms: Vec<Atom>,
}

impl Residue {
    pub fn new() -> Self {
        Residue {
            atoms: Default::default(),
        }
    }

    pub fn make(atoms: Vec<Atom>) -> Self {
        Residue { atoms }
    }
}

impl Default for Residue {
    fn default() -> Residue {
        Residue::new()
    }
}

#[derive(Serialize, Deserialize)]
pub struct Chain {
    pub residues: Vec<Residue>,
}

impl Chain {
    pub fn new() -> Self {
        Chain {
            residues: Default::default(),
        }
    }

    pub fn make(residues: Vec<Residue>) -> Self {
        Chain { residues }
    }
}

impl Default for Chain {
    fn default() -> Chain {
        Chain::new()
    }
}

#[derive(Serialize, Deserialize)]
pub struct GromacsConfig {
    pub num_steps: i64,
    pub integrator: String,
    pub seed: i64,
    pub delta_time: f64,
    pub out_steps: i64,
    pub couple_intramol: bool,
    pub stop_tolerance: f64,
}

impl GromacsConfig {
    pub fn new() -> Self {
        GromacsConfig {
            num_steps: 0,
            integrator: String::new(),
            seed: 0,
            delta_time: 0.0,
            out_steps: 0,
            couple_intramol: Default::default(),
            stop_tolerance: 0.0,
        }
    }

    pub fn make(
        num_steps: i64,
        integrator: String,
        seed: i64,
        delta_time: f64,
        out_steps: i64,
        couple_intramol: bool,
        stop_tolerance: f64,
    ) -> Self {
        GromacsConfig {
            num_steps,

            integrator,

            seed,

            delta_time,

            out_steps,

            couple_intramol,

            stop_tolerance,
        }
    }
}

impl Default for GromacsConfig {
    fn default() -> GromacsConfig {
        GromacsConfig::new()
    }
}

#[derive(Serialize, Deserialize)]
pub struct Model {
    pub chains: std::collections::HashMap<String, Chain>,
}

impl Model {
    pub fn new() -> Self {
        Model {
            chains: std::collections::HashMap::<String, Chain>::new(),
        }
    }

    pub fn make(chains: std::collections::HashMap<String, Chain>) -> Self {
        Model { chains }
    }
}

impl Default for Model {
    fn default() -> Model {
        Model::new()
    }
}

#[derive(Serialize, Deserialize)]
pub struct Structure {
    pub models: std::collections::HashMap<i64, Model>,
}

impl Structure {
    pub fn new() -> Self {
        Structure {
            models: std::collections::HashMap::<i64, Model>::new(),
        }
    }

    pub fn make(models: std::collections::HashMap<i64, Model>) -> Self {
        Structure { models }
    }
}

impl Default for Structure {
    fn default() -> Structure {
        Structure::new()
    }
}

pub mod mock {
    use super::*;

    #[derive(Clone)]
    pub struct MockBackend {
        error: Option<&'static str>,
    }

    impl MockBackend {
        pub fn new() -> MockBackend {
            MockBackend { error: None }
        }

        pub fn error(message: &'static str) -> MockBackend {
            MockBackend {
                error: Some(message),
            }
        }
    }

    #[async_trait]
    impl Backend for MockBackend {
        async fn run(&self, _: RunRequest) -> Result<RunResponse, String> {
            match self.error {
                Some(err) => Err(String::from(err)),
                None => Ok(Default::default()),
            }
        }

        async fn test(&self, _: TestRequest) -> Result<TestResponse, String> {
            match self.error {
                Some(err) => Err(String::from(err)),
                None => Ok(Default::default()),
            }
        }

        async fn visualize(&self, _: VisualizeRequest) -> Result<VisualizeResponse, String> {
            match self.error {
                Some(err) => Err(String::from(err)),
                None => Ok(Default::default()),
            }
        }
    }
}

#[cfg(test)]
mod test {
    use super::mock::*;
    use super::*;

    mod backend {
        use super::*;

        #[tokio::test]
        async fn run_ok() {
            assert!(MockBackend::new().run(RunRequest::new()).await.is_ok());
        }

        #[tokio::test]
        async fn run_error() {
            assert_eq!(
                MockBackend::error("hello from MockBackend")
                    .run(RunRequest::new())
                    .await
                    .err()
                    .unwrap(),
                "hello from MockBackend"
            );
        }

        #[tokio::test]
        async fn test_ok() {
            assert!(MockBackend::new().test(TestRequest::new()).await.is_ok());
        }

        #[tokio::test]
        async fn test_error() {
            assert_eq!(
                MockBackend::error("hello from MockBackend")
                    .test(TestRequest::new())
                    .await
                    .err()
                    .unwrap(),
                "hello from MockBackend"
            );
        }

        #[tokio::test]
        async fn visualize_ok() {
            assert!(MockBackend::new()
                .visualize(VisualizeRequest::new())
                .await
                .is_ok());
        }

        #[tokio::test]
        async fn visualize_error() {
            assert_eq!(
                MockBackend::error("hello from MockBackend")
                    .visualize(VisualizeRequest::new())
                    .await
                    .err()
                    .unwrap(),
                "hello from MockBackend"
            );
        }
    }
}
