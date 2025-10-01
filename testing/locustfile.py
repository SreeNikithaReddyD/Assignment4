from locust import FastHttpUser, task, between
import random

class ProductUser(FastHttpUser):
    wait_time = between(1, 3)  # Users wait 1-3 seconds between tasks
    created_product_ids = []

    # This task will run 9 times more often than the create_product task
    @task(9)
    def get_product(self):
        # Only try to get a product if we have created one first
        if not self.created_product_ids:
            return

        product_id = random.choice(self.created_product_ids)
        # Group all GET requests under one name in the statistics
        self.client.get(f"/products/{product_id}", name="/products/[id]")

    @task(1)
    def create_product(self):
        # Use simple, hardcoded data for the new product
        product_data = {
            "name": "Hardcoded Product",
            "price": 99.99,
            "quantity": 150
        }
        # Make the POST request and check the response
        with self.client.post("/products", json=product_data, catch_response=True) as response:
            if response.status_code == 201:
                # If creation is successful, save the new product ID
                product_id = response.json().get("id")
                if product_id:
                    self.created_product_ids.append(product_id)