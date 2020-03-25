use sal::types::GromacsConfig;
use sal::types;

struct GromacsDriver {

}

impl sal::types::Backend for GromacsDriver {
    fn run(&self, req: RunRequest) -> Result<RunResponse, String> {
        Err(String::new())
    }
}