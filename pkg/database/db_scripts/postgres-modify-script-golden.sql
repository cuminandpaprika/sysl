/*TITLE : Petstore Schema*/
/* ---------------------------------------------
Autogenerated script from sysl
--------------------------------------------- */


/*-----------------------Relation Model : RelModel-----------------------------------------------*/
ALTER TABLE Breed ALTER COLUMN breedId TYPE integer;
ALTER TABLE Breed ALTER COLUMN breedName TYPE varchar (35);
ALTER TABLE Breed ALTER COLUMN numLegs TYPE varchar (1);
CREATE TABLE Company(
  abnNumber varchar (50),
  companyName varchar (30),
  companyCountry varchar (10),
  CONSTRAINT COMPANY_PK PRIMARY KEY(abnNumber)
);
ALTER TABLE PetMedicalHistory ADD COLUMN conditionName varchar (50);
ALTER TABLE PetMedicalHistory DROP CONSTRAINT PETMEDICALHISTORY_PETID_FK;
ALTER TABLE PetMedicalHistory ALTER COLUMN petId TYPE integer;
ALTER TABLE PetMedicalHistory DROP CONSTRAINT PETMEDICALHISTORY_PK;
ALTER TABLE PetMedicalHistory ADD CONSTRAINT PETMEDICALHISTORY_PK PRIMARY KEY(conditionName,petId,reportedDate);
CREATE TABLE Department(
  deptId integer,
  deptName varchar (40),
  deptLoc varchar (50),
  abn varchar (50),
  CONSTRAINT DEPARTMENT_PK PRIMARY KEY(deptId),
  CONSTRAINT DEPARTMENT_ABN_FK FOREIGN KEY(abn) REFERENCES Company (abnNumber)
);
ALTER TABLE Pet ADD COLUMN diet varchar (45);
ALTER TABLE Pet ALTER COLUMN numLegs TYPE varchar (1);
ALTER TABLE Pet ADD COLUMN petCounter bigserial;
CREATE SEQUENCE Pet_petId_seq;
ALTER TABLE Pet ALTER COLUMN petId TYPE integer;
ALTER TABLE Pet ALTER COLUMN petId SET DEFAULT nextval('Pet_petId_seq');
ALTER SEQUENCE Pet_petId_seq OWNED BY Pet.petId;
select setval('Pet_petId_seq', coalesce(max(petId), 1)) from Pet;
ALTER TABLE Employee ADD COLUMN dept integer;
ALTER TABLE Employee ADD CONSTRAINT EMPLOYEE_DEPT_FK FOREIGN KEY(dept) REFERENCES Department (deptId);
ALTER TABLE Employee ALTER COLUMN employeeId TYPE varchar (50);
ALTER TABLE Employee ALTER COLUMN name TYPE varchar (25);
ALTER TABLE EmployeeTendsPet ALTER COLUMN employeeId TYPE varchar (50);
ALTER TABLE EmployeeTendsPet ADD CONSTRAINT EMPLOYEETENDSPET_EMPLOYEEID_FK FOREIGN KEY(employeeId) REFERENCES Employee(employeeId);
ALTER TABLE EmployeeTendsPet ALTER COLUMN petId TYPE bigint;
ALTER TABLE EmployeeTendsPet ADD CONSTRAINT EMPLOYEETENDSPET_PETID_FK FOREIGN KEY(petId) REFERENCES Pet(petId);
ALTER TABLE EmployeeTendsPet ADD COLUMN petStatus varchar (10);
ALTER TABLE EmployeeTendsPet DROP CONSTRAINT EMPLOYEETENDSPET_PK;
ALTER TABLE EmployeeTendsPet DROP COLUMN ownerShipId;
ALTER TABLE EmployeeTendsPet ADD CONSTRAINT EMPLOYEETENDSPET_PK PRIMARY KEY(employeeId,petId);
