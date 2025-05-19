1. **Activate** the virtual environment (so that `python` and `pip` point inside `.radar-oficial`):

   * **macOS/Linux**

     ```bash
     source .radar-oficial/bin/activate
     ```
2. **Verify** you’re isolated:

   ````bash
   pip list      # only shows stdlib wheels and any you’ve installed
   which pip     # points to .radar-oficial/bin/pip
   ````
3. **Install** from `requirements.txt`:

   ```bash
   pip install -r requirements.txt
   ```

   * This ensures the exact versions in your lock file are used 

4. **Add** new deps later by re-running:

   ```bash
   pip install <new-pkg>
   pip freeze > requirements.txt
   ```

   * Always regenerate the freeze file to capture updates

5. Running Your Scripts via the Virtual Environment

With `.radar-oficial` active, any invocation of `python` or your script will run inside the isolated environment:

```bash
# From project root, with .radar-oficial activated
python scripts/my_ml_script.py arg1 arg2
```